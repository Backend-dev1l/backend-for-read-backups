package interfaces

import (
	"context"
	"test-http/internal/db"
	"test-http/internal/dto"
)

type UserService interface {
	Create(ctx context.Context, request dto.CreateStatisticsRequest) (db.User, error)
	GetByID(ctx context.Context, request dto.GetUserByIDRequest) (db.User, error)
	GetByEmail(ctx context.Context, request dto.GetUserByEmailRequest) (db.User, error)
	List(ctx context.Context, request dto.ListUsersRequest) ([]db.User, error)
	Update(ctx context.Context, request dto.UpdateUserRequest) (db.User, error)
	Delete(ctx context.Context, request dto.DeleteUserRequest) error
}
