package service

import (
	"context"

	"fmt"
	"log/slog"

	"test-http/internal/db"
	"test-http/internal/lib"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserWordSetService struct {
	userWordRepo *db.Queries
	logger       *slog.Logger
}

func NewUserWordSetService(userWordRepo *db.Queries, log *slog.Logger) *UserWordSetService {
	return &UserWordSetService{
		userWordRepo: userWordRepo,
		logger:       log,
	}
}

type CreateUserWordSetParams struct {
	UserID    pgtype.UUID
	WordSetID pgtype.UUID
}

type ListUserWordSetsFilters struct {
	UserID pgtype.UUID
	Limit  int32
	Offset int32
}

func (u *UserWordSetService) Create(ctx context.Context, params CreateUserWordSetParams) (db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.Create", "creating user word set",
		slog.String("user_id", params.UserID.String()),
		slog.String("word_set_id", params.WordSetID.String()),
	)

	wordSet, err := u.userWordRepo.CreateUserWordSet(ctx, db.CreateUserWordSetParams{
		UserID:    params.UserID,
		WordSetID: params.WordSetID,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.Create", "CreateUserWordSet", "failed to create user word set", err,
			slog.String("user_id", params.UserID.String()),
			slog.String("word_set_id", params.WordSetID.String()),
		)
		return db.UserWordSet{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserWordSetService.Create", "user word set created successfully",
		slog.String("id", wordSet.ID.String()),
		slog.String("user_id", params.UserID.String()),
		slog.String("word_set_id", params.WordSetID.String()),
	)

	return wordSet, nil
}

func (u *UserWordSetService) GetByID(ctx context.Context, id pgtype.UUID) (db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.GetByID", "getting user word set by id",
		slog.String("id", id.String()),
	)

	wordSet, err := u.userWordRepo.GetUserWordSet(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.GetByID", "GetUserWordSet", "failed to get user word set by id", err,
			slog.String("id", id.String()),
		)
		return db.UserWordSet{}, InfrastructureUnexpected.Err()
	}

	return wordSet, nil
}

func (u *UserWordSetService) List(ctx context.Context, filters ListUserWordSetsFilters) ([]db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.List", "listing user word sets",
		slog.String("user_id", filters.UserID.String()),
		slog.Int("limit", int(filters.Limit)),
		slog.Int("offset", int(filters.Offset)),
	)

	wordSets, err := u.userWordRepo.ListUserWordSets(ctx, filters.UserID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.List", "ListUserWordSets", "failed to list user word sets", err,
			slog.String("user_id", filters.UserID.String()),
		)
		return nil, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserWordSetService.List", "user word sets listed successfully",
		slog.Int("count", len(wordSets)),
		slog.String("user_id", filters.UserID.String()),
	)

	return wordSets, nil
}

func (u *UserWordSetService) Update(ctx context.Context, params interface{}) (db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.Update", "update operation not implemented for user word sets")
	return db.UserWordSet{}, fmt.Errorf("update user word set not implemented")
}

func (u *UserWordSetService) Delete(ctx context.Context, id pgtype.UUID) error {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.Delete", "deleting user word set",
		slog.String("id", id.String()),
	)

	err := u.userWordRepo.DeleteUserWordSet(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.Delete", "DeleteUserWordSet", "failed to delete user word set", err,
			slog.String("id", id.String()),
		)
		return InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserWordSetService.Delete", "user word set deleted successfully",
		slog.String("id", id.String()),
	)

	return nil
}
