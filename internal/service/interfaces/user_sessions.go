package interfaces

import (
	"context"
	"test-http/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserSessionService interface {
	Create(ctx context.Context, params CreateUserSessionParams) (db.UserSession, error)
	GetByID(ctx context.Context, id pgtype.UUID) (db.UserSession, error)
	List(ctx context.Context, filters ListUserSessionsFilters) ([]db.UserSession, error)
	ListActive(ctx context.Context, userID pgtype.UUID) ([]db.UserSession, error)
	Update(ctx context.Context, params UpdateUserSessionParams) (db.UserSession, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}
