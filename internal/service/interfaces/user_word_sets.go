package interfaces

import (
	"context"
	"test-http/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserWordSetServiceInterface interface {
	Create(ctx context.Context, params CreateUserWordSetParams) (db.UserWordSet, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.UserWordSet, error)
	List(ctx context.Context, userID pgtype.UUID) ([]db.UserWordSet, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}
