package handlers

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
	"github.com/Businge931/sba-api-gateway/proto"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handler
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req proto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"details": "Invalid request"}`, http.StatusForbidden)
		return
	}

	// Validate request
	if !ValidateLoginRequest(&req) {
		http.Error(w, `{"details": "Username and password are required"}`, http.StatusForbidden)
		return
	}

	domainReq := &domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	res, err := h.authService.Login(r.Context(), domainReq)
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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req proto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"details": "Invalid request"}`, http.StatusForbidden)
		return
	}

	// Validate request
	if !ValidateRegisterRequest(&req) {
		http.Error(w, `{"details": "Username and password are required"}`, http.StatusForbidden)
		return
	}
	domainReq := &domain.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}
	res, err := h.authService.Register(r.Context(), domainReq)
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

// VerifyTokenMiddleware verifies the JWT token before allowing access to protected endpoints
func (h *AuthHandler) VerifyTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the Authorization header
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, `{"details": "Missing token"}`, http.StatusUnauthorized)
			return
		}

		// Call AuthService to verify the token
		res, err := h.authService.VerifyToken(r.Context(), token)
		if err != nil {
			handleGRPCError(w, err)
			return
		}

		// Check if the token is valid
		if !res.Success {
			http.Error(w, `{"details": "`+res.Message+`"}`, http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// ValidateLoginRequest validates the LoginRequest fields
func ValidateLoginRequest(req *proto.LoginRequest) bool {
	return req.GetUsername() != "" && req.GetPassword() != ""
}

// ValidateRegisterRequest validates the RegisterRequest fields
func ValidateRegisterRequest(req *proto.RegisterRequest) bool {
	return req.GetUsername() != "" && req.GetPassword() != ""
}

// handleGRPCError maps gRPC errors to HTTP status codes
func handleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, `{"details": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		http.Error(w, `{"details": "Invalid request"}`, http.StatusForbidden)
	case codes.Unauthenticated:
		http.Error(w, `{"details": "Unauthorized"}`, http.StatusUnauthorized)
	case codes.Internal:
		http.Error(w, `{"details": "Server error"}`, http.StatusInternalServerError)
	default:
		http.Error(w, `{"details": "Unknown error"}`, http.StatusInternalServerError)
	}
}
