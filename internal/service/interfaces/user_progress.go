package interfaces

import (
	"context"
	"test-http/internal/db"
	"test-http/internal/dto"
)

type UserProgressService interface {
	Create(ctx context.Context, request dto.CreateUserProgressRequest) (db.UserProgress, error)
	GetByID(ctx context.Context, request dto.GetUserProgressRequest) (db.UserProgress, error)
	GetByUserAndWord(ctx context.Context, request dto.GetUserProgressByUserAndWordRequest) (db.UserProgress, error)
	List(ctx context.Context, request dto.ListUserProgressRequest) ([]db.UserProgress, error)
	Update(ctx context.Context, request dto.UpdateUserProgressRequest) (db.UserProgress, error)
	Delete(ctx context.Context, request dto.DeleteUserProgressRequest) error
}
