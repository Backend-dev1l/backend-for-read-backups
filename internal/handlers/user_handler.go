package handlers

import (
	"log/slog"
	"net/http"

	"test-http/internal/dto"
	"test-http/internal/middleware"
	"test-http/internal/service"
	errorsPkg "test-http/pkg/errors_pkg"
	"test-http/pkg/fault"
	"test-http/pkg/helper"
	"test-http/pkg/uuidconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserHandler struct {
	log      *slog.Logger
	service  *service.UserService
	validate *validator.Validate
}

func NewUserHandler(log *slog.Logger, service *service.UserService) *UserHandler {
	return &UserHandler{log: log, service: service, validate: validator.New()}
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)
	log := u.log.With(slog.String("trace_id", traceID))

	log.Info("CreateUser handler started")

	defer r.Body.Close()

	var req dto.CreateUserRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		u.log.Error("DecodeJSON failed", "err", err)
		return helper.HTTPError(w, errorsPkg.DecodeFailed.Err())
	}

	if err := u.validate.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	user, err := u.service.Create(ctx, dto.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		log.Error("UserService.Create failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ContextCreatingUserMissing.Err())
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
	return nil
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)
	log := u.log.With(slog.String("trace_id", traceID))

	log.Info("GetUser handler started")

	defer r.Body.Close()

	userIDStr := chi.URLParam(r, "id")
	if userIDStr == "" {
		log.Error("missing user id in query parameters")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	pgUUID, err := uuidconv.SetPgUUID(userID)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	user, err := u.service.GetByID(ctx, dto.GetUserByIDRequest{
		ID: pgUUID,
	})
	if err != nil {
		log.Error("UserService.GetByID failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ContextGettingUserMissing.Err())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
	return nil

}

func (u *UserHandler) UserEmail(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := u.log.With(slog.String("trace_id", traceID))

	log.Info("User Email handler started")

	defer r.Body.Close()

	userEmail := chi.URLParam(r, "email")
	if userEmail == "" {
		log.Error("missing email in query parameters")
		return helper.HTTPError(w, fault.UnhandledError.Err())
	}

	user, err := u.service.GetByEmail(ctx, dto.GetUserByEmailRequest{
		Email: userEmail,
	})
	if err != nil {
		log.Error("UserService.GetByEmail failed", "email", userEmail, "err", err)
		return helper.HTTPError(w, fault.UnhandledError.Err())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
	return nil
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := u.log.With(slog.String("trace_id", traceID))
	log.Info("DeleteUser handler started")

	defer r.Body.Close()

	userIDStr := chi.URLParam(r, "id")
	if userIDStr == "" {
		log.Error("missing user id in query parameters")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	pgUUID, err := uuidconv.SetPgUUID(userID)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	if err := u.service.Delete(ctx, dto.DeleteUserRequest{
		ID: pgUUID,
	}); err != nil {
		log.Error("UserService.Delete failed", "err", err)
		return helper.HTTPError(w, fault.UnhandledError.Err())
	}

	render.Status(r, http.StatusOK)
	return nil
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	traceID := middleware.GetTraceID(ctx)

	log := u.log.With(slog.String("trace_id", traceID))
	log.Info("UpdateUser handler started")

	defer r.Body.Close()

	var req dto.UpdateUserRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("DecodeJSON failed", "err", err)
		return helper.HTTPError(w, errorsPkg.DecodeFailed.Err())
	}

	if err := u.validate.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	userIDStr := chi.URLParam(r, "id")
	if userIDStr == "" {
		log.Error("missing user id in query parameters")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	pgUUID, err := uuidconv.SetPgUUID(userID)
	if err != nil {
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	user, err := u.service.Update(ctx, dto.UpdateUserRequest{
		ID: pgUUID,
	})
	if err != nil {
		log.Error("UserService.Update failed", "err", err)
		return helper.HTTPError(w, fault.UnhandledError.Err())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)

	return nil
}
