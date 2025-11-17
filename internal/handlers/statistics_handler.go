package handlers

import (
	"log/slog"
	"net/http"
	"test-http/internal/dto"
	"test-http/internal/middleware"
	"test-http/internal/service"
	errorsPkg "test-http/pkg/errors_pkg"

	"test-http/pkg/helper"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	log.Info("CreateStatistics handler called")

	defer r.Body.Close()

	var req dto.CreateStatisticsRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("DecodeJSON failed", "err", err)
		return helper.HTTPError(w, errorsPkg.DecodeFailed.Err())
	}

	log.Info("Decoded request", "user_id", req.UserID.String(), "total_words_learned", req.TotalWordsLearned, "accuracy", req.Accuracy, "total_time", req.TotalTime)

	if err := s.validate.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	if !req.UserID.Valid {
		log.Error("invalid user_id")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	statistics, err := s.service.Create(ctx, req)
	if err != nil {
		log.Error("UserStatisticsService.Create failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ContextCreatingUserStatisticsMissing.Err())
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, statistics)
	return nil
}

func (s *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := s.logger.With(slog.String("trace_id", traceID))

	log.Info("GetStatistics handler called")

	defer r.Body.Close()

	userIDStr := chi.URLParam(r, "user_id")
	if userIDStr == "" {
		log.Error("missing user_id in query parameters")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	if _, err := uuid.Parse(userIDStr); err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	userID, err := helper.ToUUID(userIDStr)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	statistics, err := s.service.GetByID(ctx, dto.GetStatisticsRequest{
		UserID: userID,
	})
	if err != nil {
		log.Error("UserStatisticsService.GetByID failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ContextGettingUserMissing.Err())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, statistics)
	return nil
}

func (s *StatisticsHandler) UpdateStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := s.logger.With(slog.String("trace_id", traceID))

	log.Info("UpdateStatistics handler called")

	defer r.Body.Close()

	var req dto.UpdateStatisticsRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("DecodeJSON failed", "err", err)
		return helper.HTTPError(w, errorsPkg.DecodeFailed.Err())
	}

	if err := s.validate.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	if !req.UserID.Valid {
		log.Error("invalid user_id")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	statistics, err := s.service.Update(ctx, req)
	if err != nil {
		log.Error("UserStatisticsService.Update failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ContextUpdatingUserStatisticsMissing.Err())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, statistics)
	return nil
}

func (s *StatisticsHandler) DeleteStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := s.logger.With(slog.String("trace_id", traceID))

	log.Info("DeleteStatistics handler called")

	defer r.Body.Close()

	userIDStr := chi.URLParam(r, "id")
	if userIDStr == "" {
		log.Error("missing user_id in query parameters")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	if _, err := uuid.Parse(userIDStr); err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	userID, err := helper.ToUUID(userIDStr)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	if err := s.service.Delete(ctx, dto.DeleteStatisticsRequest{
		UserID: userID,
	}); err != nil {
		log.Error("UserStatisticsService.Delete failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ContextDeletingUserStatisticsMissing.New())
	}

	render.Status(r, http.StatusOK)
	return nil
}
