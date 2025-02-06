package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/proto"
)

type OddsClient struct {
	client proto.OddsServiceClient
}

func NewOddsClient(conn *grpc.ClientConn) *OddsClient {
	return &OddsClient{client: proto.NewOddsServiceClient(conn)}
}

func (c *OddsClient) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	protoReq := &proto.CreateOddsRequest{
		League:          req.League,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: float32(req.HomeTeamWinOdds),
		AwayTeamWinOdds: float32(req.AwayTeamWinOdds),
		DrawOdds:        float32(req.DrawOdds),
		GameDate:        req.GameDate,
	}

	res, err := c.client.CreateOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &domain.CreateOddsResponse{
		Success: res.Success,
		Message: res.Message,
		Details: res.Details,
	}, nil
}

func (c *OddsClient) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	protoReq := &proto.ReadOddsRequest{
		League: req.League,
		Date:   req.Date,
	}

	res, err := c.client.ReadOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// Map the gRPC response to the domain response
	odds := make([]domain.CreateOddsRequest, len(res.Odds))
	for i, protoOdds := range res.Odds {
		odds[i] = domain.CreateOddsRequest{
			League:          protoOdds.League,
			HomeTeam:        protoOdds.HomeTeam,
			AwayTeam:        protoOdds.AwayTeam,
			HomeTeamWinOdds: float64(protoOdds.HomeTeamWinOdds),
			AwayTeamWinOdds: float64(protoOdds.AwayTeamWinOdds),
			DrawOdds:        float64(protoOdds.DrawOdds),
			GameDate:        protoOdds.GameDate,
		}
	}

	return &domain.ReadOddsResponse{
		Odds:    odds,
		Details: res.Details,
	}, nil
}
func (c *OddsClient) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	protoReq := &proto.UpdateOddsRequest{
		League:          req.League,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: float32(req.HomeTeamWinOdds),
		AwayTeamWinOdds: float32(req.AwayTeamWinOdds),
		DrawOdds:        float32(req.DrawOdds),
		GameDate:        req.GameDate,
	}

	res, err := c.client.UpdateOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &domain.UpdateOddsResponse{
		Success: res.Success,
		Message: res.Message,
		Details: res.Details,
	}, nil
}

func (c *OddsClient) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	protoReq := &proto.DeleteOddsRequest{
		League:   req.League,
		HomeTeam: req.HomeTeam,
		AwayTeam: req.AwayTeam,
		GameDate: req.GameDate,
	}

	res, err := c.client.DeleteOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &domain.DeleteOddsResponse{
		Success: res.Success,
		Message: res.Message,
		Details: res.Details,
	}, nil
}
