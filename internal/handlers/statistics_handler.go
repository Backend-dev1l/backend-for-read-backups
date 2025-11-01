package handlers

import (
	"log/slog"
	"net/http"
	"test-http/internal/db"
	"test-http/internal/middleware"
	"test-http/internal/service"
	errorsPkg "test-http/pkg/errors_pkg"

	"test-http/pkg/helper"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type StatisticsHandler struct {
	logger   *slog.Logger
	validate *validator.Validate
	service  *service.UserStatisticsService
}

func NewStatisticsHandler(service *service.UserStatisticsService, validate *validator.Validate, logger *slog.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		logger:   logger,
		validate: validate,
		service:  service,
	}
}

func (s *StatisticsHandler) CreateStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := s.logger.With(slog.String("trace_id", traceID))

	log.Info("Statistics handler called")

	var req db.UserStatistic

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("DecodeJSON failed:", err)
		return helper.HTTPError(w, errorsPkg.DecodeFailed.Err())
	}

	if err := s.validate.Struct(req); err != nil {
		log.Error("validation failed:", err)
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	statistics, err := s.service.Create(ctx, service.CreateUserStatisticsParams{
		UserID:            req.UserID,
		TotalWordsLearned: req.TotalWordsLearned,
		Accuracy:          req.Accuracy,
		TotalTime:         req.TotalTime,
	})
	if err != nil {
		log.Error("UserStatisticsService.Create failed:", err)
		return helper.HTTPError(w, errorsPkg.ContextCreatingUserStatisticsMissing.Err())
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, statistics)
	return nil

}
