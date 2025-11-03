package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	routes "test-http/cmd/api"
	"test-http/internal/config"

	pool "test-http/pkg/db"
	"test-http/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.InitLogger(cfg)
	if err != nil {
		panic(err)
	}

	log.Info("Logger succesfully initialized", slog.String("level", cfg.Logger.Level))

	dbPool, err := pool.NewPool(cfg, log)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	r := chi.NewRouter()

	routes.RegisterRoutes(r, dbPool, cfg, log)

	srv := &http.Server{
		Addr:         cfg.Address(),
		IdleTimeout:  cfg.TimeOuts.IdleTimeout,
		ReadTimeout:  cfg.TimeOuts.ReadTimeout,
		WriteTimeout: cfg.TimeOuts.WriteTimeout,
		Handler:      r,
	}

	log.Info("Starting server", slog.String("address", cfg.Address()))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", slog.String("error", err.Error()))

		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-done
	log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.TimeOuts.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", slog.String("error", err.Error()))
	} else {
		log.Info("Graceful shutdown completed")
	}

}
