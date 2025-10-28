package service

import (
	"context"

	"log/slog"

	"test-http/internal/db"
	"test-http/internal/lib"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserStatisticsService struct {
	userStatistRepo *db.Queries
	logger          *slog.Logger
}

func NewUserStatisticsService(userStatistRepo *db.Queries, log *slog.Logger) *UserStatisticsService {
	return &UserStatisticsService{
		userStatistRepo: userStatistRepo,
		logger:          log,
	}
}

type CreateUserStatisticsParams struct {
	UserID            pgtype.UUID
	TotalWordsLearned int32
	Accuracy          pgtype.Numeric
	TotalTime         int32
}

type UpdateUserStatisticsParams struct {
	UserID            pgtype.UUID
	TotalWordsLearned int32
	Accuracy          pgtype.Numeric
	TotalTime         int32
}

func (u *UserStatisticsService) Create(ctx context.Context, params CreateUserStatisticsParams) (db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.Create", "creating user statistics",
		slog.String("user_id", params.UserID.String()),
		slog.Int("total_words_learned", int(params.TotalWordsLearned)),
		slog.Int("total_time", int(params.TotalTime)),
	)

	stats, err := u.userStatistRepo.CreateUserStatistics(ctx, db.CreateUserStatisticsParams{
		UserID:            params.UserID,
		TotalWordsLearned: params.TotalWordsLearned,
		Accuracy:          params.Accuracy,
		TotalTime:         params.TotalTime,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Create", "CreateUserStatistics", "failed to create user statistics", err,
			slog.String("user_id", params.UserID.String()),
		)
		return db.UserStatistic{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.Create", "user statistics created successfully",
		slog.String("user_id", stats.UserID.String()),
		slog.Int("total_words_learned", int(stats.TotalWordsLearned)),
	)

	return stats, nil
}

func (u *UserStatisticsService) GetByID(ctx context.Context, userID pgtype.UUID) (db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.GetByID", "getting user statistics by user id",
		slog.String("user_id", userID.String()),
	)

	stats, err := u.userStatistRepo.GetUserStatistics(ctx, userID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.GetByID", "GetUserStatistics", "failed to get user statistics by user id", err,
			slog.String("user_id", userID.String()),
		)
		return db.UserStatistic{}, InfrastructureUnexpected.Err()
	}

	return stats, nil
}

func (u *UserStatisticsService) List(ctx context.Context, filters UpdateUserStatisticsParams) ([]db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.List", "list operation not implemented for user statistics")

	statistics, err := u.userStatistRepo.ListActiveSessions(ctx)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.List", "failed to list user statistics", err,
	   slog.String("user_id", filters.UserID.String())),
		 return nil, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.List", "successfully listed users",
    slog.String("user_id", filters.UserID.String()),
	  slog.Int("total_words_learned", int(params.TotalWordsLearned))
	  slog.Int("total_time", int(params.TotalTime))),

	return statistics, nil
}

func (u *UserStatisticsService) Update(ctx context.Context, params UpdateUserStatisticsParams) (db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.Update", "updating user statistics",
		slog.String("user_id", params.UserID.String()),
		slog.Int("total_words_learned", int(params.TotalWordsLearned)),
		slog.Int("total_time", int(params.TotalTime)),
	)

	stats, err := u.userStatistRepo.UpdateUserStatistics(ctx, db.UpdateUserStatisticsParams{
		UserID:            params.UserID,
		TotalWordsLearned: params.TotalWordsLearned,
		Accuracy:          params.Accuracy,
		TotalTime:         params.TotalTime,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Update", "UpdateUserStatistics", "failed to update user statistics", err,
			slog.String("user_id", params.UserID.String()),
		)
		return db.UserStatistic{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.Update", "user statistics updated successfully",
		slog.String("user_id", stats.UserID.String()),
		slog.Int("total_words_learned", int(stats.TotalWordsLearned)),
	)

	return stats, nil
}

func (u *UserStatisticsService) Delete(ctx context.Context, userID pgtype.UUID) error {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.Delete", "deleting user statistics",
		slog.String("user_id", userID.String()),
	)

	err := u.userStatistRepo.DeleteUserStatistics(ctx, userID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Delete", "DeleteUserStatistics", "failed to delete user statistics", err,
			slog.String("user_id", userID.String()),
		)
		return InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.Delete", "user statistics deleted successfully",
		slog.String("user_id", userID.String()),
	)

	return nil
}
