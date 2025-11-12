package service

import (
	"context"

	"log/slog"

	"test-http/internal/db"
	"test-http/internal/dto"
	"test-http/internal/lib"
	errorsPkg "test-http/pkg/errors_pkg"

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

func (u *UserSessionService) Create(ctx context.Context, request dto.CreateUserSessionRequest) (db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.Create", "creating user session",
		slog.String("user_id", request.UserID.String()),
		slog.String("status", request.Status),
	)

	session, err := u.userSessionRepo.CreateUserSession(ctx, db.CreateUserSessionParams{
		UserID: request.UserID,
		Status: request.Status,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.Create", "CreateUserSession", "failed to create user session", err,
			slog.String("user_id", request.UserID.String()),
			slog.String("status", request.Status),
		)
		return db.UserSession{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.Create", "user session created successfully",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", request.UserID.String()),
		slog.String("status", request.Status),
	)

	return session, nil
}

func (u *UserSessionService) GetByID(ctx context.Context, request dto.GetUserSessionRequest) (db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.GetByID", "getting user session by id",
		slog.String("session_id", request.ID.String()),
	)

	session, err := u.userSessionRepo.GetUserSession(ctx, request.ID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionRepository.GetByID", "GetUserSession", "failed to get user session by id", err,
			slog.String("session_id", request.ID.String()),
		)
		return db.UserSession{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	return session, nil
}

func (u *UserSessionService) List(ctx context.Context, request dto.ListUserSessionsRequest) ([]db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionRepository.List", "listing user sessions",
		slog.String("user_id", request.UserID.String()),
		slog.Int("limit", int(request.Limit)),
		slog.Int("offset", int(request.Offset)),
	)

	sessions, err := u.userSessionRepo.ListUserSessions(ctx, db.ListUserSessionsParams{
		UserID: request.UserID,
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.List", "ListUserSessions", "failed to list user sessions", err,
			slog.String("user_id", request.UserID.String()),
		)
		return nil, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.List", "user sessions listed successfully",
		slog.Int("count", len(sessions)),
		slog.String("user_id", request.UserID.String()),
	)

	return sessions, nil
}

func (u *UserSessionService) ListActive(ctx context.Context) ([]db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.ListActive", "listing active sessions")

	sessions, err := u.userSessionRepo.ListActiveSessions(ctx, db.ListActiveSessionsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.ListActive", "ListActiveSessions", "failed to list active sessions", err)
		return nil, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.ListActive", "active sessions listed successfully",
		slog.Int("count", len(sessions)),
	)

	return sessions, nil
}

func (u *UserSessionService) Update(ctx context.Context, request dto.UpdateUserSessionRequest) (db.UserSession, error) {
	lib.LogDebug(ctx, u.logger, "UserSessionService.Update", "updating user session",
		slog.String("session_id", request.ID.String()),
		slog.String("status", request.Status),
	)

	session, err := u.userSessionRepo.UpdateUserSession(ctx, db.UpdateUserSessionParams{
		ID:      request.ID,
		Status:  request.Status,
		EndedAt: request.EndedAt,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserSessionService.Update", "UpdateUserSession", "failed to update user session", err,
			slog.String("session_id", request.ID.String()),
			slog.String("status", request.Status),
		)
		return db.UserSession{}, errorsPkg.InfrastructureUnexpected.Err()
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
		return errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserSessionService.Delete", "user session deleted successfully",
		slog.String("session_id", id.String()),
	)

	return nil
}
