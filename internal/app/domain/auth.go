package domain

import (
	"context"
)


type AuthClient interface {
	// Login authenticates a user and returns a token.
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	// Register creates a new user account.
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)

	// VerifyToken validates a JWT token and returns its status.
	VerifyToken(ctx context.Context, token string) (*VerifyTokenResponse, error)
}

// LoginRequest represents the request payload for the Login method.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response payload for the Login method.
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// RegisterRequest represents the request payload for the Register method.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse represents the response payload for the Register method.
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// VerifyTokenRequest represents the request payload for the VerifyToken method.
type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// VerifyTokenResponse represents the response payload for the VerifyToken method.
type VerifyTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}