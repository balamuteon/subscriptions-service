// @title Subscriptions Service API
// @version 1.0
// @description REST API for managing user subscriptions and calculating totals.
// @BasePath /api/v1
// @schemes http
// @host localhost:8080
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "subscription_service/docs"
	"subscription_service/internal/config"
	"subscription_service/internal/httpapi"
	subscriptionHandler "subscription_service/internal/httpapi"
	subscriptionRepo "subscription_service/internal/repository/subscription"
	"subscription_service/internal/server"
	subscriptionService "subscription_service/internal/service/subscription"
	"subscription_service/pkg/logger"
	"subscription_service/pkg/postgres"
)

func main() {
	// 1. Init configuration
	cfg, err := config.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Init logger
	log := logger.New(cfg.Server.LogLevel)

	// 3. Init db
	db, err := postgres.NewConnection(context.Background(), cfg.Database.DSN())
	if err != nil {
		log.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 4. Init deps (repository, service, HTTP handlers)
	repo := subscriptionRepo.New(db)
	service := subscriptionService.New(repo)
	handler := subscriptionHandler.NewSubscriptionHandler(log, service)

	// 5. Init HTTP router and server
	router := httpapi.NewHandler(log.With("component", "http"), handler)
	srv := server.New(cfg.Server, router)

	errCh := make(chan error, 1)
	go func() {
		log.Info("http server started", "addr", srv.Addr)
		// 6. Run HTTP server
		errCh <- srv.ListenAndServe()
	}()

	// 7. Listen shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Info("shutdown signal received", "signal", sig.String())
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped")
}