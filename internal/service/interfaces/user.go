package interfaces

import (
	"context"
	"test-http/internal/db"
	"test-http/internal/dto"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService interface {
	Create(ctx context.Context, request dto.CreateStatisticsRequest) (db.User, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	List(ctx context.Context, filters ListUsersFilters) ([]db.User, error)
	Update(ctx context.Context, params UpdateUserParams) (db.User, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}
