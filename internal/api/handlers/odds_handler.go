package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
	"github.com/Businge931/sba-api-gateway/proto"
)

type OddsHandler struct {
	oddsService service.OddsService
}

func NewOddsHandler(oddsService service.OddsService) *OddsHandler {
	return &OddsHandler{oddsService: oddsService}
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
	serviceCall func(context.Context, interface{}) (interface{}, error),
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

	// Call the service method
	res, err := serviceCall(r.Context(), req)
	if err != nil {
		handleGRPCError(w, err)
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
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.CreateOddsRequest)
		domainReq := &domain.CreateOddsRequest{
			League:          protoReq.League,
			GameDate:        protoReq.GameDate,
			HomeTeam:        protoReq.HomeTeam,
			AwayTeam:        protoReq.AwayTeam,
			HomeTeamWinOdds: float64(protoReq.HomeTeamWinOdds),
			AwayTeamWinOdds: float64(protoReq.AwayTeamWinOdds),
			DrawOdds:        float64(protoReq.DrawOdds),
		}
		return h.oddsService.CreateOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.CreateOddsRequest{})
}

// ReadOdds handles the reading of odds
func (h *OddsHandler) ReadOdds(w http.ResponseWriter, r *http.Request) {
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.ReadOddsRequest)
		domainReq := &domain.ReadOddsRequest{
			League: protoReq.League,
			Date:   protoReq.Date,
		}
		return h.oddsService.ReadOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.ReadOddsRequest{})
}

// UpdateOdds handles the updating of odds
func (h *OddsHandler) UpdateOdds(w http.ResponseWriter, r *http.Request) {
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.UpdateOddsRequest)
		domainReq := &domain.UpdateOddsRequest{
			League:          protoReq.League,
			GameDate:        protoReq.GameDate,
			HomeTeam:        protoReq.HomeTeam,
			AwayTeam:        protoReq.AwayTeam,
			HomeTeamWinOdds: float64(protoReq.HomeTeamWinOdds),
			AwayTeamWinOdds: float64(protoReq.AwayTeamWinOdds),
			DrawOdds:        float64(protoReq.DrawOdds),
		}
		return h.oddsService.UpdateOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.UpdateOddsRequest{})
}

// DeleteOdds handles the deletion of odds
func (h *OddsHandler) DeleteOdds(w http.ResponseWriter, r *http.Request) {
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.DeleteOddsRequest)
		domainReq := &domain.DeleteOddsRequest{
			League:   protoReq.League,
			GameDate: protoReq.GameDate,
		}
		return h.oddsService.DeleteOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.DeleteOddsRequest{})
}