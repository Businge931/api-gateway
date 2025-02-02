package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/Businge931/sba-api-gateway/proto"
)

type OddsHandler struct {
	client proto.OddsServiceClient
}

func NewOddsHandler(conn *grpc.ClientConn) *OddsHandler {
	return &OddsHandler{client: proto.NewOddsServiceClient(conn)}
}

func validateLeagueAndDate(league, date string) (string, bool) {
	if league != "english premier league" {
		return `{"details": "League must be 'english premier league'"}`, false
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return `{"details": "Invalid date format. Use YYYY-MM-DD"}`, false
	}

	return "", true
}

// handleValidationError writes a validation error response
func handleValidationError(w http.ResponseWriter, errMsg string) {
	http.Error(w, errMsg, http.StatusForbidden)
}

// handleOddsRequest is a helper function to handle common logic for odds requests
func (h *OddsHandler) handleOddsRequest(
	w http.ResponseWriter,
	r *http.Request,
	grpcCall func(context.Context, interface{}) (interface{}, error),
	req interface{},
) {
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handleValidationError(w, `{"details": "Invalid request"}`)
		return
	}

	// Validate league and date
	var league, date string
	switch r := req.(type) {
	case *proto.CreateOddsRequest:
		league, date = r.GetLeague(), r.GetGameDate()
	case *proto.ReadOddsRequest:
		league, date = r.GetLeague(), r.GetDate()
	case *proto.UpdateOddsRequest:
		league, date = r.GetLeague(), r.GetGameDate()
	case *proto.DeleteOddsRequest:
		league, date = r.GetLeague(), r.GetGameDate()
	}

	if errMsg, valid := validateLeagueAndDate(league, date); !valid {
		handleValidationError(w, errMsg)
		return
	}

	// Call the gRPC method
	res, err := grpcCall(r.Context(), req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Check if the operation was successful
	var success bool
	var details string
	switch r := res.(type) {
	case *proto.CreateOddsResponse:
		success, details = r.GetSuccess(), r.GetDetails()
	case *proto.ReadOddsResponse:
		success, details = true, "" // Read operation always considered successful if we reach here
	case *proto.UpdateOddsResponse:
		success, details = r.GetSuccess(), r.GetDetails()
	case *proto.DeleteOddsResponse:
		success, details = r.GetSuccess(), r.GetDetails()
	}

	if !success {
		handleValidationError(w, `{"details": "`+details+`"}`)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

// CreateOdds handles the creation of odds
func (h *OddsHandler) CreateOdds(w http.ResponseWriter, r *http.Request) {
	h.handleOddsRequest(w, r, func(ctx context.Context, req interface{}) (interface{}, error) {
		return h.client.CreateOdds(ctx, req.(*proto.CreateOddsRequest))
	}, &proto.CreateOddsRequest{})
}

// ReadOdds handles the updating of odds
func (h *OddsHandler) ReadOdds(w http.ResponseWriter, r *http.Request) {
	h.handleOddsRequest(w, r, func(ctx context.Context, req interface{}) (interface{}, error) {
		return h.client.ReadOdds(ctx, req.(*proto.ReadOddsRequest))
	}, &proto.ReadOddsRequest{})
}

// UpdateOdds handles the updating of odds
func (h *OddsHandler) UpdateOdds(w http.ResponseWriter, r *http.Request) {
	h.handleOddsRequest(w, r, func(ctx context.Context, req interface{}) (interface{}, error) {
		return h.client.UpdateOdds(ctx, req.(*proto.UpdateOddsRequest))
	}, &proto.UpdateOddsRequest{})
}

// DeleteOdds handles the deletion of odds
func (h *OddsHandler) DeleteOdds(w http.ResponseWriter, r *http.Request) {
	h.handleOddsRequest(w, r, func(ctx context.Context, req interface{}) (interface{}, error) {
		return h.client.DeleteOdds(ctx, req.(*proto.DeleteOddsRequest))
	}, &proto.DeleteOddsRequest{})
}
