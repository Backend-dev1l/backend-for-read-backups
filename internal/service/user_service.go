package service

import (
	"context"
	"strings"

	"log/slog"
	"test-http/internal/db"
	"test-http/internal/dto"
	"test-http/internal/lib"
	errorsPkg "test-http/pkg/errors_pkg"
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

func (u *UserService) Create(ctx context.Context, request dto.CreateUserRequest) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.Create", "creating user",
		slog.String("username", request.Username),
		slog.String("email", request.Email),
	)

	user, err := u.userRepo.CreateUser(ctx, db.CreateUserParams{
		Username: request.Username,
		Email:    request.Email,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.Create", "CreateUser", "failed to create user", err,
			slog.String("username", request.Username),
			slog.String("email", request.Email),
		)
		return db.User{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.Create", "user created successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
	)

	return user, nil
}

func (u *UserService) GetByID(ctx context.Context, request dto.GetUserByIDRequest) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.GetByID", "getting user by id",
		slog.String("user_id", request.ID.String()),
	)

	user, err := u.userRepo.GetUser(ctx, request.ID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.GetByID", "GetUser", "failed to get user by id", err,
			slog.String("user_id", request.ID.String()),
		)
		return db.User{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	return user, nil
}

func (u *UserService) GetByEmail(ctx context.Context, request dto.GetUserByEmailRequest) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.GetByEmail", "getting user by email",
		slog.String("email", request.Email),
	)

	if !strings.Contains(request.Email, "@") {
		lib.LogError(ctx, u.logger, "UserService.GetByEmail", "GetUserByEmail", "invalid email format", nil,
			slog.String("email", request.Email))
		return db.User{}, errorsPkg.ValidationError.Err()
	}

	user, err := u.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.GetByEmail", "GetUserByEmail", "failed to get user by email", err,
			slog.String("email", request.Email),
		)
		return db.User{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	return user, nil
}

func (u *UserService) List(ctx context.Context, request dto.ListUsersRequest) ([]db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.List", "listing users",
		slog.Int("limit", int(request.Limit)),
		slog.Int("offset", int(request.Offset)),
	)

	users, err := u.userRepo.ListUsers(ctx, db.ListUsersParams{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.List", "ListUsers", "failed to list users", err,
			slog.Int("limit", int(request.Limit)),
			slog.Int("offset", int(request.Offset)),
		)
		return nil, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.List", "users listed successfully",
		slog.Int("count", len(users)),
		slog.Int("limit", int(request.Limit)),
		slog.Int("offset", int(request.Offset)),
	)

	return users, nil
}

func (u *UserService) Update(ctx context.Context, request dto.UpdateUserRequest) (db.User, error) {
	lib.LogDebug(ctx, u.logger, "UserService.Update", "updating user",
		slog.String("user_id", request.ID.String()),
		slog.String("username", request.Username),
		slog.String("email", request.Email),
	)

	user, err := u.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		ID:       request.ID,
		Username: request.Username,
		Email:    request.Email,
	})
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.Update", "UpdateUser", "failed to update user", err,
			slog.String("user_id", request.ID.String()),
			slog.String("username", request.Username),
			slog.String("email", request.Email),
		)
		return db.User{}, errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.Update", "user updated successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
	)

	return user, nil
}

func (u *UserService) Delete(ctx context.Context, request dto.DeleteUserRequest) error {
	lib.LogDebug(ctx, u.logger, "UserService.Delete", "deleting user",
		slog.String("user_id", request.ID.String()),
	)

	err := u.userRepo.DeleteUser(ctx, request.ID)
	if err != nil {
		lib.LogError(ctx, u.logger, "UserService.Delete", "DeleteUser", "failed to delete user", err,
			slog.String("user_id", request.ID.String()),
		)
		return errorsPkg.InfrastructureUnexpected.Err()
	}

	lib.LogInfo(ctx, u.logger, "UserService.Delete", "user deleted successfully",
		slog.String("user_id", request.ID.String()),
	)

	return nil
}
