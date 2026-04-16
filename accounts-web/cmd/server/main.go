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

	"github.com/solardome/gamepulse-platform/accounts-web/internal/config"
	"github.com/solardome/gamepulse-platform/accounts-web/internal/graphql"
	webhttp "github.com/solardome/gamepulse-platform/accounts-web/internal/http"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	graphQLClient := graphql.New(cfg.BackendGraphQLURL, cfg.RequestTimeout)
	serverHandler, err := webhttp.NewServer(logger, graphQLClient)
	if err != nil {
		logger.Error("build accounts-web server failed", "err", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           serverHandler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("accounts-web started", "addr", cfg.HTTPAddr, "backend_graphql_url", cfg.BackendGraphQLURL)
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Error("accounts-web serve failed", "err", serveErr)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("accounts-web shutting down")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("accounts-web shutdown failed", "err", err)
		os.Exit(1)
	}
}
