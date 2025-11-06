package service

import (
	"context"

	"log/slog"

	"test-http/internal/db"
	"test-http/internal/lib"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserSessionService struct {
	userSessionRepo db.UserSessionRepo
	logger          *slog.Logger
}

func NewUserSessionService(userSessionRepo db.UserSessionRepo, log *slog.Logger) *UserSessionService {
	return &UserSessionService{
		userSessionRepo: userSessionRepo,
		logger:          log,
	}
}

type CreateUserSessionParams struct {
	UserID pgtype.UUID
	Status string
}

type UpdateUserSessionParams struct {
	ID      pgtype.UUID
	Status  string
	EndedAt pgtype.Timestamptz
}

type ListUserSessionsFilters struct {
	UserID pgtype.UUID
	Limit  int32
	Offset int32
}

func (u *UserSessionService) Create(ctx context.Context, params CreateUserSessionParams) (db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.Create", "creating user session",
		slog.String("user_id", params.UserID.String()),
		slog.String("status", params.Status),
	)

	session, err := u.userSessionRepo.CreateUserSession(ctx, db.CreateUserSessionParams{
		UserID: params.UserID,
		Status: params.Status,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.Create", "CreateUserSession", "failed to create user session", err,
			slog.String("user_id", params.UserID.String()),
			slog.String("status", params.Status),
		)
		return db.UserSession{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.Create", "user session created successfully",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", params.UserID.String()),
		slog.String("status", params.Status),
	)

	return session, nil
}

func (u *UserSessionService) GetByID(ctx context.Context, id pgtype.UUID) (db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.GetByID", "getting user session by id",
		slog.String("session_id", id.String()),
	)

	session, err := u.userSessionRepo.GetUserSession(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionRepository.GetByID", "GetUserSession", "failed to get user session by id", err,
			slog.String("session_id", id.String()),
		)
		return db.UserSession{}, InfrastructureUnexpected.Err()
	}

	return session, nil
}

func (u *UserSessionService) List(ctx context.Context, filters ListUserSessionsFilters) ([]db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionRepository.List", "listing user sessions",
		slog.String("user_id", filters.UserID.String()),
		slog.Int("limit", int(filters.Limit)),
		slog.Int("offset", int(filters.Offset)),
	)

	sessions, err := u.userSessionRepo.ListUserSessions(ctx, filters.UserID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.List", "ListUserSessions", "failed to list user sessions", err,
			slog.String("user_id", filters.UserID.String()),
		)
		return nil, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.List", "user sessions listed successfully",
		slog.Int("count", len(sessions)),
		slog.String("user_id", filters.UserID.String()),
	)

	return sessions, nil
}

func (u *UserSessionService) ListActive(ctx context.Context) ([]db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.ListActive", "listing active sessions")

	sessions, err := u.userSessionRepo.ListActiveSessions(ctx)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.ListActive", "ListActiveSessions", "failed to list active sessions", err)
		return nil, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.ListActive", "active sessions listed successfully",
		slog.Int("count", len(sessions)),
	)

	return sessions, nil
}

func (u *UserSessionService) Update(ctx context.Context, params UpdateUserSessionParams) (db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.Update", "updating user session",
		slog.String("session_id", params.ID.String()),
		slog.String("status", params.Status),
	)

	session, err := u.userSessionRepo.UpdateUserSession(ctx, db.UpdateUserSessionParams{
		ID:      params.ID,
		Status:  params.Status,
		EndedAt: params.EndedAt,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.Update", "UpdateUserSession", "failed to update user session", err,
			slog.String("session_id", params.ID.String()),
			slog.String("status", params.Status),
		)
		return db.UserSession{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.Update", "user session updated successfully",
		slog.String("session_id", session.ID.String()),
		slog.String("status", session.Status),
	)

	return session, nil
}

func (u *UserSessionService) Delete(ctx context.Context, id pgtype.UUID) error {
	lib.LogDebug(ctx, u.logger, "UserSessionRepository.Delete", "deleting user session",
		slog.String("session_id", id.String()),
	)

	err := u.userSessionRepo.DeleteUserSession(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.Delete", "DeleteUserSession", "failed to delete user session", err,
			slog.String("session_id", id.String()),
		)
		return InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.Delete", "user session deleted successfully",
		slog.String("session_id", id.String()),
	)

	return nil
}
