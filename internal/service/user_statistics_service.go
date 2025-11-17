package service

import (
	"context"

	"log/slog"

	"test-http/internal/db"
	"test-http/internal/dto"
	"test-http/internal/lib"
	errorsPkg "test-http/pkg/errors_pkg"
	"test-http/pkg/helper"
)

type UserStatisticsService struct {
	userStatistRepo db.UserStatisticsRepo
	logger          *slog.Logger
}

func NewUserStatisticsService(userStatistRepo db.UserStatisticsRepo, log *slog.Logger) *UserStatisticsService {
	return &UserStatisticsService{
		userStatistRepo: userStatistRepo,
		logger:          log,
	}
}

func (u *UserStatisticsService) Create(ctx context.Context, request dto.CreateStatisticsRequest) (db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.Create", "creating user statistics",
		slog.String("user_id", request.UserID.String()),
		slog.Int("total_words_learned", int(request.TotalWordsLearned)),
		slog.Int("total_time", int(request.TotalTime)),
	)

	accuracy, err := helper.ToNumeric(request.Accuracy)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Create", "ToNumeric", "failed to convert accuracy", err,
			slog.String("user_id", request.UserID.String()),
		)
		return db.UserStatistic{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	stats, err := u.userStatistRepo.CreateUserStatistics(ctx, db.CreateUserStatisticsParams{
		UserID:            request.UserID,
		TotalWordsLearned: request.TotalWordsLearned,
		Accuracy:          accuracy,
		TotalTime:         request.TotalTime,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Create", "CreateUserStatistics", "failed to create user statistics", err,
			slog.String("user_id", request.UserID.String()),
		)
		return db.UserStatistic{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.Create", "user statistics created successfully",
		slog.String("user_id", stats.UserID.String()),
		slog.Int("total_words_learned", int(stats.TotalWordsLearned)),
	)

	return stats, nil
}

func (u *UserStatisticsService) GetByID(ctx context.Context, request dto.GetStatisticsRequest) (db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.GetByID", "getting user statistics by user id",
		slog.String("user_id", request.UserID.String()),
	)

	stats, err := u.userStatistRepo.GetUserStatistics(ctx, request.UserID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.GetByID", "GetUserStatistics", "failed to get user statistics by user id", err,
			slog.String("user_id", request.UserID.String()),
		)
		return db.UserStatistic{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	return stats, nil
}

func (u *UserStatisticsService) List(ctx context.Context, request dto.ListStatisticsRequest) ([]db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.List", "listing user statistics",
		slog.Int("limit", int(request.Limit)),
		slog.Int("offset", int(request.Offset)),
	)

	stats, err := u.userStatistRepo.ListUserStatistics(ctx, db.ListUserStatisticsParams{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.List", "ListUserStatistics", "failed to list user statistics", err)
		return nil, errorsPkg.InfrastructureUnexpected.Err()
	}
	lib.LogInfo(ctx, u.logger, "UserStatisticsService.List", "list operation completed",
		slog.Int("count", len(stats)),
	)

	return stats, nil
}

func (u *UserStatisticsService) Update(ctx context.Context, request dto.UpdateStatisticsRequest) (db.UserStatistic, error) {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.Update", "updating user statistics",
		slog.String("user_id", request.UserID.String()),
		slog.Int("total_words_learned", int(request.TotalWordsLearned)),
		slog.Int("total_time", int(request.TotalTime)),
	)

	accuracy, err := helper.ToNumeric(request.Accuracy)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Update", "ToNumeric", "failed to convert accuracy", err,
			slog.String("user_id", request.UserID.String()),
		)
		return db.UserStatistic{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	stats, err := u.userStatistRepo.UpdateUserStatistics(ctx, db.UpdateUserStatisticsParams{
		UserID:            request.UserID,
		TotalWordsLearned: request.TotalWordsLearned,
		Accuracy:          accuracy,
		TotalTime:         request.TotalTime,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Update", "UpdateUserStatistics", "failed to update user statistics", err,
			slog.String("user_id", request.UserID.String()),
		)
		return db.UserStatistic{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.Update", "user statistics updated successfully",
		slog.String("user_id", stats.UserID.String()),
		slog.Int("total_words_learned", int(stats.TotalWordsLearned)),
	)

	return stats, nil
}

func (u *UserStatisticsService) Delete(ctx context.Context, request dto.DeleteStatisticsRequest) error {
	lib.LogDebug(ctx, u.logger, "UserStatisticsService.Delete", "deleting user statistics",
		slog.String("user_id", request.UserID.String()),
	)

	err := u.userStatistRepo.DeleteUserStatistics(ctx, request.UserID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserStatisticsService.Delete", "DeleteUserStatistics", "failed to delete user statistics", err,
			slog.String("user_id", request.UserID.String()),
		)
		return errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserStatisticsService.Delete", "user statistics deleted successfully",
		slog.String("user_id", request.UserID.String()),
	)

	return nil
}
