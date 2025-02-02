package handlers

import (
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

func handleValidationError(w http.ResponseWriter, errMsg string) {
	http.Error(w, errMsg, http.StatusForbidden)
}

func (handler *OddsHandler) CreateOdds(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleValidationError(w, `{"details": "Invalid request"}`)
		return
	}

	if errMsg, valid := validateLeagueAndDate(req.GetLeague(), req.GetGameDate()); !valid {
		handleValidationError(w, errMsg)
		return
	}

	res, err := handler.client.CreateOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	if !res.GetSuccess() {
		handleValidationError(w, `{"details": "`+res.GetDetails()+`"}`)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (handler *OddsHandler) ReadOdds(w http.ResponseWriter, r *http.Request) {
	var req proto.ReadOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleValidationError(w, `{"details": "Invalid request"}`)
		return
	}

	if _, err := time.Parse("2006-01-02", req.GetDate()); err != nil {
		handleValidationError(w, `{"details": "Invalid date format. Use YYYY-MM-DD"}`)
		return
	}

	res, err := handler.client.ReadOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (handler *OddsHandler) UpdateOdds(w http.ResponseWriter, r *http.Request) {
	var req proto.UpdateOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleValidationError(w, `{"details": "Invalid request"}`)
		return
	}

	if errMsg, valid := validateLeagueAndDate(req.GetLeague(), req.GetGameDate()); !valid {
		handleValidationError(w, errMsg)
		return
	}

	res, err := handler.client.UpdateOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	if !res.GetSuccess() {
		handleValidationError(w, `{"details": "`+res.GetDetails()+`"}`)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (handler *OddsHandler) DeleteOdds(w http.ResponseWriter, r *http.Request) {
	var req proto.DeleteOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleValidationError(w, `{"details": "Invalid request"}`)
		return
	}

	if errMsg, valid := validateLeagueAndDate(req.GetLeague(), req.GetGameDate()); !valid {
		handleValidationError(w, errMsg)
		return
	}

	res, err := handler.client.DeleteOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	if !res.GetSuccess() {
		handleValidationError(w, `{"details": "`+res.GetDetails()+`"}`)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
