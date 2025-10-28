package service

import (
	"context"

	"log/slog"

	"test-http/internal/db"
	"test-http/internal/lib"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserProgressService struct {
	userProgressRepo *db.Queries
	logger           *slog.Logger
}

func NewUserProgressService(userProgresRepo *db.Queries, log *slog.Logger) *UserProgressService {
	return &UserProgressService{
		userProgressRepo: userProgresRepo,
		logger:           log,
	}
}

type CreateUserProgressParams struct {
	UserID         pgtype.UUID
	WordID         pgtype.UUID
	CorrectCount   int32
	IncorrectCount int32
}

type UpdateUserProgressParams struct {
	ID             pgtype.UUID
	CorrectCount   int32
	IncorrectCount int32
}

type ListUserProgressFilters struct {
	UserID pgtype.UUID
	Limit  int32
	Offset int32
}

func (u *UserProgressService) Create(ctx context.Context, params CreateUserProgressParams) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.Create", "creating user progress",
		slog.String("user_id", params.UserID.String()),
		slog.String("word_id", params.WordID.String()),
		slog.Int("correct_count", int(params.CorrectCount)),
		slog.Int("incorrect_count", int(params.IncorrectCount)),
	)

	progress, err := u.userProgressRepo.CreateUserProgress(ctx, db.CreateUserProgressParams{
		UserID:         params.UserID,
		WordID:         params.WordID,
		CorrectCount:   params.CorrectCount,
		IncorrectCount: params.IncorrectCount,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.Create", "CreateUserProgress", "failed to create user progress", err,
			slog.String("user_id", params.UserID.String()),
			slog.String("word_id", params.WordID.String()),
		)
		return db.UserProgress{}, InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.Create", "user progress created successfully",
		slog.String("progress_id", progress.ID.String()),
	)

	return progress, nil
}

func (u *UserProgressService) GetByID(ctx context.Context, id pgtype.UUID) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.GetByID", "getting user progress by id",
		slog.String("progress_id", id.String()),
	)

	progress, err := u.userProgressRepo.GetUserProgress(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.GetByID", "GetUserProgress", "failed to get user progress by id", err,
			slog.String("progress_id", id.String()),
		)
		return db.UserProgress{}, InfrastructureUnexpected.Err()
	}

	return progress, nil
}

func (u *UserProgressService) GetByUserAndWord(ctx context.Context, userID, wordID pgtype.UUID) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.GetByUserAndWord", "getting user progress by user and word",
		slog.String("user_id", userID.String()),
		slog.String("word_id", wordID.String()),
	)

	progress, err := u.userProgressRepo.GetUserProgressByUserAndWord(ctx, db.GetUserProgressByUserAndWordParams{
		UserID: userID,
		WordID: wordID,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.GetByUserAndWord", "GetUserProgressByUserAndWord", "failed to get user progress by user and word", err,
			slog.String("user_id", userID.String()),
			slog.String("word_id", wordID.String()),
		)
		return db.UserProgress{}, InfrastructureUnexpected.Err()
	}

	return progress, nil
}

func (u *UserProgressService) List(ctx context.Context, filters ListUserProgressFilters) ([]db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.List", "listing user progress",
		slog.String("user_id", filters.UserID.String()),
		slog.Int("limit", int(filters.Limit)),
		slog.Int("offset", int(filters.Offset)),
	)

	progressList, err := u.userProgressRepo.ListUserProgress(ctx, filters.UserID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.List", "ListUserProgress", "failed to list user progress", err,
			slog.String("user_id", filters.UserID.String()),
		)
		return nil, InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.List", "user progress listed successfully",
		slog.Int("count", len(progressList)),
	)

	return progressList, nil
}

func (u *UserProgressService) Update(ctx context.Context, params UpdateUserProgressParams) (db.UserProgress, error) {
	lib.LogDebug(ctx, u.logger, "UserProgressService.Update", "updating user progress",
		slog.String("progress_id", params.ID.String()),
		slog.Int("correct_count", int(params.CorrectCount)),
		slog.Int("incorrect_count", int(params.IncorrectCount)),
	)

	progress, err := u.userProgressRepo.UpdateUserProgress(ctx, db.UpdateUserProgressParams{
		ID:             params.ID,
		CorrectCount:   params.CorrectCount,
		IncorrectCount: params.IncorrectCount,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.Update", "UpdateUserProgress", "failed to update user progress", err,
			slog.String("progress_id", params.ID.String()),
		)
		return db.UserProgress{}, InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.Update", "user progress updated successfully",
		slog.String("progress_id", progress.ID.String()),
	)

	return progress, nil
}

func (u *UserProgressService) Delete(ctx context.Context, id pgtype.UUID) error {
	lib.LogDebug(ctx, u.logger, "UserProgressService.Delete", "deleting user progress",
		slog.String("progress_id", id.String()),
	)

	err := u.userProgressRepo.DeleteUserProgress(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserProgressService.Delete", "DeleteUserProgress", "failed to delete user progress", err,
			slog.String("progress_id", id.String()),
		)
		return InfrastructureUnexpected.Err()
	}

	lib.LogDebug(ctx, u.logger, "UserProgressService.Delete", "user progress deleted successfully",
		slog.String("progress_id", id.String()),
	)

	return nil
}
