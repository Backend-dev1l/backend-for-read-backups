package interfaces

import (
	"context"
	"test-http/internal/db"
	"test-http/internal/dto"
)

type UserSessionService interface {
	Create(ctx context.Context, request dto.CreateUserSessionRequest) (db.UserSession, error)
	GetByID(ctx context.Context, request dto.GetUserSessionRequest) (db.UserSession, error)
	List(ctx context.Context, request dto.ListUserSessionsRequest) ([]db.UserSession, error)
	ListActive(ctx context.Context, request dto.ListActiveUserSessionsRequest) ([]db.UserSession, error)
	Update(ctx context.Context, request dto.UpdateUserSessionRequest) (db.UserSession, error)
	Delete(ctx context.Context, request dto.DeleteUserSessionRequest) error
}
