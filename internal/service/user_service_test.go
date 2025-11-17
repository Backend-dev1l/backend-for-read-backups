package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"test-http/internal/db"
	"test-http/internal/db/mocks"
	"test-http/internal/dto"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	// minimal no-op logger for tests
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserService(mockRepo, logger)

	params := db.CreateUserParams{Username: "alice", Email: "alice@example.com"}
	want := db.User{Username: "alice", Email: "alice@example.com"}

	mockRepo.EXPECT().CreateUser(gomock.Any(), params).Return(want, nil)

	got, err := svc.Create(context.Background(), dto.CreateUserRequest{Username: params.Username, Email: params.Email})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != want.Username || got.Email != want.Email {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserService_GetByEmail_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserService(mockRepo, logger)

	_, err := svc.GetByEmail(context.Background(), dto.GetUserByEmailRequest{Email: "bad-email"})
	if err == nil {
		t.Fatalf("expected validation error for invalid email, got nil")
	}
}

func TestUserService_GetByEmail_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserService(mockRepo, logger)

	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "nope@example.com").Return(db.User{}, errors.New("not found"))

	_, err := svc.GetByEmail(context.Background(), dto.GetUserByEmailRequest{Email: "nope@example.com"})
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserService(mockRepo, logger)

	users := []db.User{{Username: "a"}, {Username: "b"}}
	mockRepo.EXPECT().ListUsers(gomock.Any(), db.ListUsersParams{Limit: 10, Offset: 0}).Return(users, nil)

	got, err := svc.List(context.Background(), dto.ListUsersRequest{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(users) {
		t.Fatalf("got %d users want %d", len(got), len(users))
	}
}

func TestUserService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserService(mockRepo, logger)

	// Use zero UUID for simplicity
	var id pgtype.UUID

	up := db.UpdateUserParams{Username: "bob", Email: "bob@ex.com", ID: id}
	want := db.User{Username: up.Username, Email: up.Email}

	// We expect UpdateUser with any context and the params passed
	mockRepo.EXPECT().UpdateUser(gomock.Any(), up).Return(want, nil)

	got, err := svc.Update(context.Background(), dto.UpdateUserRequest{ID: id, Username: up.Username, Email: up.Email})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != want.Username || got.Email != want.Email {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserService(mockRepo, logger)

	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	var id pgtype.UUID
	if err := id.Scan(uuidStr); err != nil {
		t.Fatalf("failed to convert uuid")
	}

	mockRepo.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)

	err := svc.Delete(context.Background(), dto.DeleteUserRequest{ID: id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
