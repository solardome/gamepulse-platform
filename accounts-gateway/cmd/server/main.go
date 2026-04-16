package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/solardome/gamepulse-platform/accounts-gateway/graph"
	"github.com/solardome/gamepulse-platform/accounts-gateway/internal/account"
	"github.com/solardome/gamepulse-platform/accounts-gateway/internal/config"
	httpserver "github.com/solardome/gamepulse-platform/accounts-gateway/internal/http"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	accountClient, err := account.New(cfg.AccountGRPCAddr, cfg.RequestTimeout)
	if err != nil {
		logger.Error("create account client failed", "err", err)
		os.Exit(1)
	}
	defer accountClient.Close()

	router := httpserver.NewRouter(logger, &graph.Resolver{
		AccountClient: accountClient,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("accounts-gateway started", "addr", cfg.HTTPAddr, "account_grpc_addr", cfg.AccountGRPCAddr)
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Error("accounts-gateway http serve failed", "err", serveErr)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("accounts-gateway shutting down")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("accounts-gateway shutdown failed", "err", err)
		os.Exit(1)
	}
}
