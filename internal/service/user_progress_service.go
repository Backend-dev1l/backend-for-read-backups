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

type UserProgressService struct {
	userProgressRepo db.UserProgressRepo
	logger           *slog.Logger
}

func NewUserProgressService(userProgressRepo db.UserProgressRepo, log *slog.Logger) *UserProgressService {
	return &UserProgressService{
		userProgressRepo: userProgressRepo,
		logger:           log,
	}
}

type ListUserProgressFilters struct {
	UserID pgtype.UUID
	Limit  int32
	Offset int32
}

func (u *UserProgressService) Create(ctx context.Context, request dto.CreateUserProgressRequest) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.Create", "creating user progress",
		slog.String("user_id", request.UserID.String()),
		slog.String("word_id", request.WordID.String()),
		slog.Int("correct_count", int(request.CorrectCount)),
		slog.Int("incorrect_count", int(request.IncorrectCount)),
	)

	progress, err := u.userProgressRepo.CreateUserProgress(ctx, db.CreateUserProgressParams{
		UserID:         request.UserID,
		WordID:         request.WordID,
		CorrectCount:   request.CorrectCount,
		IncorrectCount: request.IncorrectCount,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.Create", "CreateUserProgress", "failed to create user progress", err,
			slog.String("user_id", request.UserID.String()),
			slog.String("word_id", request.WordID.String()),
		)
		return db.UserProgress{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.Create", "user progress created successfully",
		slog.String("progress_id", progress.ID.String()),
	)

	return progress, nil
}

func (u *UserProgressService) GetByID(ctx context.Context, request dto.GetUserProgressRequest) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.GetByID", "getting user progress by id",
		slog.String("progress_id", request.ID.String()),
	)

	progress, err := u.userProgressRepo.GetUserProgress(ctx, request.ID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.GetByID", "GetUserProgress", "failed to get user progress by id", err,
			slog.String("progress_id", request.ID.String()),
		)
		return db.UserProgress{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	return progress, nil
}

func (u *UserProgressService) GetByUserAndWord(ctx context.Context, request dto.GetUserProgressByUserAndWordRequest) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.GetByUserAndWord", "getting user progress by user and word",
		slog.String("user_id", request.UserID.String()),
		slog.String("word_id", request.WordID.String()),
	)

	progress, err := u.userProgressRepo.GetUserProgressByUserAndWord(ctx, db.GetUserProgressByUserAndWordParams{
		UserID: request.UserID,
		WordID: request.WordID,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.GetByUserAndWord", "GetUserProgressByUserAndWord", "failed to get user progress by user and word", err,
			slog.String("user_id", request.UserID.String()),
			slog.String("word_id", request.WordID.String()),
		)
		return db.UserProgress{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	return progress, nil
}

func (u *UserProgressService) List(ctx context.Context, request dto.ListUserProgressRequest) ([]db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.List", "listing user progress",
		slog.String("user_id", request.UserID.String()),
		slog.Int("limit", int(request.Limit)),
		slog.Int("offset", int(request.Offset)),
	)

	progressList, err := u.userProgressRepo.ListUserProgress(ctx, request.UserID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.List", "ListUserProgress", "failed to list user progress", err,
			slog.String("user_id", request.UserID.String()),
		)
		return nil, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.List", "user progress listed successfully",
		slog.Int("count", len(progressList)),
	)

	return progressList, nil
}

func (u *UserProgressService) Update(ctx context.Context, request dto.UpdateUserProgressRequest) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.Update", "updating user progress",
		slog.String("progress_id", request.ID.String()),
		slog.Int("correct_count", int(request.CorrectCount)),
		slog.Int("incorrect_count", int(request.IncorrectCount)),
	)

	progress, err := u.userProgressRepo.UpdateUserProgress(ctx, db.UpdateUserProgressParams{
		ID:             request.ID,
		CorrectCount:   request.CorrectCount,
		IncorrectCount: request.IncorrectCount,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.Update", "UpdateUserProgress", "failed to update user progress", err,
			slog.String("progress_id", request.ID.String()),
		)
		return db.UserProgress{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.Update", "user progress updated successfully",
		slog.String("progress_id", progress.ID.String()),
	)

	return progress, nil
}

func (u *UserProgressService) Delete(ctx context.Context, request dto.DeleteUserProgressRequest) error {
	lib.LogDebug(ctx, u.logger, "UserProgressService.Delete", "deleting user progress",
		slog.String("progress_id", request.ID.String()),
	)

	err := u.userProgressRepo.DeleteUserProgress(ctx, request.ID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.Delete", "DeleteUserProgress", "failed to delete user progress", err,
			slog.String("progress_id", request.ID.String()),
		)
		return errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.Delete", "user progress deleted successfully",
		slog.String("progress_id", request.ID.String()),
	)

	return nil
}
