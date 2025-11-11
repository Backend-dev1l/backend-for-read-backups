package service

import (
	"context"
	"log/slog"

	"test-http/internal/db"
	"test-http/internal/dto"
	"test-http/internal/lib"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserWordSetService struct {
	userWordRepo db.UserWordSetRepo
	logger       *slog.Logger
}

func NewUserWordSetService(userWordRepo db.UserWordSetRepo, log *slog.Logger) *UserWordSetService {
	return &UserWordSetService{
		userWordRepo: userWordRepo,
		logger:       log,
	}
}

func (u *UserWordSetService) Create(ctx context.Context, request dto.CreateUserWordSetRequest) (db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.Create", "creating user word set",
		slog.String("user_id", request.UserID.String()),
		slog.String("word_set_id", request.WordSetID.String()),
	)

	wordSet, err := u.userWordRepo.CreateUserWordSet(ctx, db.CreateUserWordSetParams{
		UserID:    request.UserID,
		WordSetID: request.WordSetID,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.Create", "CreateUserWordSet", "failed to create user word set", err,
			slog.String("user_id", request.UserID.String()),
			slog.String("word_set_id", request.WordSetID.String()),
		)
		return db.UserWordSet{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserWordSetService.Create", "user word set created successfully",
		slog.String("id", wordSet.ID.String()),
		slog.String("user_id", request.UserID.String()),
		slog.String("word_set_id", request.WordSetID.String()),
	)

	return wordSet, nil
}

func (u *UserWordSetService) GetByID(ctx context.Context, request dto.GetUserWordSetRequest) (db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.GetByID", "getting user word set by id",
		slog.String("id", request.ID.String()),
	)

	wordSet, err := u.userWordRepo.GetUserWordSet(ctx, request.ID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.GetByID", "GetUserWordSet", "failed to get user word set by id", err,
			slog.String("id", request.ID.String()),
		)
		return db.UserWordSet{}, InfrastructureUnexpected.Err()
	}

	return wordSet, nil
}

func (u *UserWordSetService) List(ctx context.Context, request dto.ListUserWordSetsRequest) ([]db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.List", "listing user word sets",
		slog.String("user_id", request.UserID.String()),
		slog.Int("limit", int(request.Limit)),
		slog.Int("offset", int(request.Offset)),
	)

	wordSets, err := u.userWordRepo.ListUserWordSets(ctx, db.ListUserWordSetsParams{
		UserID: request.UserID,
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.List", "ListUserWordSets", "failed to list user word sets", err,
			slog.String("user_id", request.UserID.String()),
		)
		return nil, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserWordSetService.List", "user word sets listed successfully",
		slog.Int("count", len(wordSets)),
		slog.String("user_id", request.UserID.String()),
	)

	return wordSets, nil
}

func (u *UserWordSetService) Update(ctx context.Context, request dto.UpdateUserWordSetRequest) (db.UserWordSet, error) {
	lib.LogDebug(ctx, u.logger, "UserWordSetService.Update", "updating user word set",
		slog.String("id", request.ID.String()),
		slog.String("word_set_id", request.WordSetID.String()),
	)

	wordSet, err := u.userWordRepo.UpdateUserWordSet(ctx, db.UpdateUserWordSetParams{
		ID:        request.ID,
		WordSetID: request.WordSetID,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserWordSetService.Update", "UpdateUserWordSet", "failed to update user word set", err,
			slog.String("id", request.ID.String()),
			slog.String("word_set_id", request.WordSetID.String()),
		)
		return db.UserWordSet{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserWordSetService.Update", "user word set updated successfully",
		slog.String("id", wordSet.ID.String()),
		slog.String("word_set_id", wordSet.WordSetID.String()),
	)

	return wordSet, nil
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
