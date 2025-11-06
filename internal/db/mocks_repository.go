package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepo defines the subset of methods from Queries used by services.
// This interface exists to allow mocking in tests.
type UserRepo interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, id pgtype.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	CountUsers(ctx context.Context) (int64, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

// UserSessionRepo defines methods from Queries used by session services.
type UserSessionRepo interface {
	CreateUserSession(ctx context.Context, arg CreateUserSessionParams) (UserSession, error)
	GetUserSession(ctx context.Context, id pgtype.UUID) (UserSession, error)
	ListActiveSessions(ctx context.Context, arg ListActiveSessionsParams) ([]UserSession, error)
	ListUserSessions(ctx context.Context, arg ListUserSessionsParams) ([]UserSession, error)
	UpdateUserSession(ctx context.Context, arg UpdateUserSessionParams) (UserSession, error)
	DeleteUserSession(ctx context.Context, id pgtype.UUID) error
}

// UserProgressRepo defines methods from Queries used by progress services.
type UserProgressRepo interface {
	CreateUserProgress(ctx context.Context, arg CreateUserProgressParams) (UserProgress, error)
	GetUserProgress(ctx context.Context, id pgtype.UUID) (UserProgress, error)
	GetUserProgressByUserAndWord(ctx context.Context, arg GetUserProgressByUserAndWordParams) (UserProgress, error)
	ListUserProgress(ctx context.Context, userID pgtype.UUID) ([]UserProgress, error)
	UpdateUserProgress(ctx context.Context, arg UpdateUserProgressParams) (UserProgress, error)
	DeleteUserProgress(ctx context.Context, id pgtype.UUID) error
}

// UserStatisticsRepo defines methods from Queries used by statistics service.
type UserStatisticsRepo interface {
	CreateUserStatistics(ctx context.Context, arg CreateUserStatisticsParams) (UserStatistic, error)
	GetUserStatistics(ctx context.Context, userID pgtype.UUID) (UserStatistic, error)
	UpdateUserStatistics(ctx context.Context, arg UpdateUserStatisticsParams) (UserStatistic, error)
	DeleteUserStatistics(ctx context.Context, userID pgtype.UUID) error
}

// UserWordSetRepo defines methods from Queries used by user word set service.
type UserWordSetRepo interface {
	CreateUserWordSet(ctx context.Context, arg CreateUserWordSetParams) (UserWordSet, error)
	GetUserWordSet(ctx context.Context, id pgtype.UUID) (UserWordSet, error)
	ListUserWordSets(ctx context.Context, userID pgtype.UUID) ([]UserWordSet, error)
	DeleteUserWordSet(ctx context.Context, id pgtype.UUID) error
}
