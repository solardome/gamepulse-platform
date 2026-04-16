package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/solardome/gamepulse-platform/accounts-service/internal/config"
	grpcserver "github.com/solardome/gamepulse-platform/accounts-service/internal/server"
	"github.com/solardome/gamepulse-platform/accounts-service/internal/service"
	accountv1 "github.com/solardome/gamepulse-platform/gen/account/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Error("listen failed", "err", err)
		os.Exit(1)
	}

	svc := service.New()
	server := grpc.NewServer()
	accountv1.RegisterAccountServiceServer(server, grpcserver.New(logger, svc))

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	healthv1.RegisterHealthServer(server, healthServer)

	reflection.Register(server)

	go func() {
		logger.Info("accounts-service started", "addr", cfg.GRPCAddr)
		if serveErr := server.Serve(listener); serveErr != nil {
			logger.Error("grpc serve failed", "err", serveErr)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("accounts-service shutting down")
	server.GracefulStop()
}
