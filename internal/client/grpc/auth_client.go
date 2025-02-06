package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/proto"
)

type AuthClient struct {
	client proto.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{client: proto.NewAuthServiceClient(conn)}
}

func (c *AuthClient) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	protoReq := &proto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	res, err := c.client.Login(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.LoginResponse{
		Token: res.Token,
	}, nil
}

func (c *AuthClient) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	protoReq := &proto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}
	res, err := c.client.Register(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.RegisterResponse{
		Message: res.Message,
	}, nil
}

func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	res, err := c.client.VerifyToken(ctx, &proto.VerifyTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	return &domain.VerifyTokenResponse{
		Success: res.Success,
		Message: res.Message,
	}, nil
}
