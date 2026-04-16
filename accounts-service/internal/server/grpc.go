package server

import (
	"context"
	"log/slog"

	"github.com/solardome/gamepulse-platform/accounts-service/internal/service"
	accountv1 "github.com/solardome/gamepulse-platform/gen/account/v1"
)

type GRPCServer struct {
	accountv1.UnimplementedAccountServiceServer

	logger  *slog.Logger
	service *service.AccountService
}

func New(logger *slog.Logger, svc *service.AccountService) *GRPCServer {
	return &GRPCServer{
		logger:  logger,
		service: svc,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, req *accountv1.PingRequest) (*accountv1.PingResponse, error) {
	s.logger.InfoContext(ctx, "accounts-service ping request", "origin", req.GetOrigin())

	return s.service.Ping(req.GetOrigin()), nil
}
