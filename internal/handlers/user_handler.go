package handlers

import (
	"log/slog"
	"net/http"

	"test-http/internal/db"

	"test-http/internal/service"
	errorsPkg "test-http/pkg/errors_pkg"
	"test-http/pkg/helper"
	"test-http/pkg/logger"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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

	traceID, _ := ctx.Value(logger.TraceIDKey).(string)
	log := u.log.With(
		slog.String("trace_id", traceID),
	)

	log.Info("CreateUser handler started")

	var req db.User
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("DecodeJSON failed:", err)
		return helper.HTTPError(ctx, w, r, errorsPkg.DecodeFailed.Err())
	}

	if err := u.validate.Struct(req); err != nil {
		log.Error("validation failed:", err)
		return helper.HTTPError(ctx, w, r, errorsPkg.ValidationError.Err())
	}

	user, err := u.service.Create(ctx, service.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		log.Error("UserService.Create failed:", err)
		return helper.HTTPError(ctx, w, r, errorsPkg.ContextCreatingUserMissing.Err())
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
	return nil
}
