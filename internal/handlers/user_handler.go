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
	pgtype "github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
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

	var req db.User
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		u.log.Error("DecodeJSON failed:", err)
		return helper.HTTPError(w, errorsPkg.DecodeFailed.Err())
	}

	if err := u.validate.Struct(req); err != nil {
		log.Error("validation failed:", err)
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	user, err := u.service.Create(ctx, service.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		log.Error("UserService.Create failed:", err)
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

	userID := r.URL.Query().Get("id")
	if userID == "" {
		log.Error("missing user id in query parameters")
		return helper.HTTPError(w, errorsPkg.ValidationError.Err())
	}

	parsedUUID, err := uuid.FromString(userID)
	if err != nil {
		log.Error("invalid user id format:", err)
		return helper.HTTPError(w, errorsPkg.UUIDParsingFailed.Err())
	}

	userUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	user, err := u.service.GetByID(ctx, userUUID)
	if err != nil {
		log.Error("UserService.GetByID failed:", err)
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

	userEmail := r.URL.Query().Get("email")
	if userEmail != "" {
		log.Error("missing email in query parameters")
		return helper.HTTPError(w, errorsPkg.GettingEmailError.Err())
	}

	user, err := u.service.GetByEmail(ctx, userEmail)
	if err != nil {
		log.Error("UserService.GetByEmail failed:", err)
		return helper.HTTPError(w, errorsPkg.ContextGettingUserMissing.Err())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
	return nil
}
