package service

import (
	"context"
	"strings"

	"log/slog"

	"test-http/internal/db"
	"test-http/internal/lib"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	userRepo db.UserRepo
	logger   *slog.Logger
}

func NewUserService(userRepo db.UserRepo, log *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   log,
	}
}

type ListUsersFilters struct {
	Limit  int32
	Offset int32
}

type UpdateUserParams struct {
	ID       pgtype.UUID
	Username string
	Email    string
}

type CreateUserParams struct {
	Username string
	Email    string
}

func (u *UserService) Create(ctx context.Context, params CreateUserParams) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.Create", "creating user",
		slog.String("username", params.Username),
		slog.String("email", params.Email),
	)

	user, err := u.userRepo.CreateUser(ctx, db.CreateUserParams{
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.Create", "CreateUser", "failed to create user", err,
			slog.String("username", params.Username),
			slog.String("email", params.Email),
		)
		return db.User{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.Create", "user created successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
	)

	return user, nil
}

func (u *UserService) GetByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.GetByID", "getting user by id",
		slog.String("user_id", id.String()),
	)

	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.GetByID", "GetUser", "failed to get user by id", err,
			slog.String("user_id", id.String()),
		)
		return db.User{}, InfrastructureUnexpected.Err()
	}

	return user, nil
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.GetByEmail", "getting user by email",
		slog.String("email", email),
	)

	if !strings.Contains(email, "@") {
		lib.LogError(ctx, u.logger, "UserService.GetByEmail", "GetUserByEmail", "invalid email format", nil,
			slog.String("email", email))
		return db.User{}, ValidationFailed.Err()
	}

	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.GetByEmail", "GetUserByEmail", "failed to get user by email", err,
			slog.String("email", email),
		)
		return db.User{}, InfrastructureUnexpected.Err()
	}

	return user, nil
}

func (u *UserService) List(ctx context.Context, filters ListUsersFilters) ([]db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.List", "listing users",
		slog.Int("limit", int(filters.Limit)),
		slog.Int("offset", int(filters.Offset)),
	)

	users, err := u.userRepo.ListUsers(ctx, db.ListUsersParams{
		Limit:  filters.Limit,
		Offset: filters.Offset,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.List", "ListUsers", "failed to list users", err,
			slog.Int("limit", int(filters.Limit)),
			slog.Int("offset", int(filters.Offset)),
		)
		return nil, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.List", "users listed successfully",
		slog.Int("count", len(users)),
		slog.Int("limit", int(filters.Limit)),
		slog.Int("offset", int(filters.Offset)),
	)

	return users, nil
}

func (u *UserService) Update(ctx context.Context, params UpdateUserParams) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.Update", "updating user",
		slog.String("user_id", params.ID.String()),
		slog.String("username", params.Username),
		slog.String("email", params.Email),
	)

	user, err := u.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		ID:       params.ID,
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.Update", "UpdateUser", "failed to update user", err,
			slog.String("user_id", params.ID.String()),
			slog.String("username", params.Username),
			slog.String("email", params.Email),
		)
		return db.User{}, InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.Update", "user updated successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
	)

	return user, nil
}

func (u *UserService) Delete(ctx context.Context, id pgtype.UUID) error {
	lib.LogDebug(ctx, u.logger, "UserService.Delete", "deleting user",
		slog.String("user_id", id.String()),
	)

	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.Delete", "DeleteUser", "failed to delete user", err,
			slog.String("user_id", id.String()),
		)
		return InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.Delete", "user deleted successfully",
		slog.String("user_id", id.String()),
	)

	return nil
}
