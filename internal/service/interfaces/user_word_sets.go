package interfaces

import (
	"context"
	"test-http/internal/db"
	"test-http/internal/dto"
)

type UserWordSetServiceInterface interface {
	Create(ctx context.Context, request dto.CreateUserWordSetRequest) (db.UserWordSet, error)
	GetByID(ctx context.Context, request dto.GetUserWordSetRequest) (db.UserWordSet, error)
	List(ctx context.Context, request dto.ListUserWordSetsRequest) ([]db.UserWordSet, error)
	Delete(ctx context.Context, request dto.DeleteUserWordSetRequest) error
}
