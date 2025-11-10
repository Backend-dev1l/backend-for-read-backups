package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"test-http/internal/config"
	"test-http/internal/db"
	"test-http/internal/handlers"
	"test-http/internal/middleware"
	"test-http/internal/service"
)

func RegisterRoutes(r chi.Router, dbPool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) {
	validate := validator.New()
	
	userRepo := db.New(dbPool)
	userService := service.NewUserService(userRepo, logger)
	userHandler := handlers.NewUserHandler(logger, userService)

	userStatisticsService := service.NewUserStatisticsService(userRepo, logger)
	statisticsHandler := handlers.NewStatisticsHandler(userStatisticsService, validate, logger)

	readyHandler := handlers.NewReadyHandler(dbPool, cfg, logger)

	r.Use(middleware.TraceID)
	r.Use(middleware.Recover(logger))
	r.Use(middleware.RequestLogger(logger))

	r.Route("/api/v1", func(r chi.Router) {
		// --- Users ---
		r.Route("/users", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { _ = userHandler.CreateUser(w, r) })
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) { _ = userHandler.GetUser(w, r) })
			r.Get("/email/{email}", func(w http.ResponseWriter, r *http.Request) { _ = userHandler.UserEmail(w, r) })
			r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) { _ = userHandler.UpdateUser(w, r) })
			r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { _ = userHandler.DeleteUser(w, r) })
		})

		// --- Statistics ---
		r.Route("/statistics", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { _ = statisticsHandler.CreateStatistics(w, r) })
		})
	})

	// --- Health Check ---
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) { _ = readyHandler.ReadyzHandler(w, r) })
}
