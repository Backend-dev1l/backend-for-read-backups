package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	db "test-http/internal/db"
	mocks "test-http/internal/db/mocks"
	"test-http/internal/dto"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserSessionService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var uid pgtype.UUID
	params := dto.CreateUserSessionRequest{UserID: uid, Status: "active"}
	want := db.UserSession{UserID: uid, Status: "active"}

	mockRepo.EXPECT().CreateUserSession(gomock.Any(), db.CreateUserSessionParams{UserID: uid, Status: "active"}).Return(want, nil)

	got, err := svc.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != want.UserID || got.Status != want.Status {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserSessionService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var uid pgtype.UUID
	params := dto.CreateUserSessionRequest{UserID: uid, Status: "active"}

	mockRepo.EXPECT().CreateUserSession(gomock.Any(), gomock.Any()).Return(db.UserSession{}, errors.New("db error"))

	_, err := svc.Create(context.Background(), params)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserSessionService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var id pgtype.UUID
	want := db.UserSession{ID: id, Status: "active"}
	mockRepo.EXPECT().GetUserSession(gomock.Any(), id).Return(want, nil)

	got, err := svc.GetByID(context.Background(), dto.GetUserSessionRequest{ID: id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID || got.Status != want.Status {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserSessionService_GetByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().GetUserSession(gomock.Any(), id).Return(db.UserSession{}, errors.New("not found"))

	_, err := svc.GetByID(context.Background(), dto.GetUserSessionRequest{ID: id})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserSessionService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var uid pgtype.UUID
	users := []db.UserSession{{UserID: uid}, {UserID: uid}}
	mockRepo.EXPECT().ListUserSessions(gomock.Any(), db.ListUserSessionsParams{UserID: uid, Limit: 10, Offset: 0}).Return(users, nil)

	got, err := svc.List(context.Background(), dto.ListUserSessionsRequest{UserID: uid, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(users) {
		t.Fatalf("got %d want %d", len(got), len(users))
	}
}

func TestUserSessionService_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var uid pgtype.UUID
	mockRepo.EXPECT().ListUserSessions(gomock.Any(), db.ListUserSessionsParams{UserID: uid, Limit: 1, Offset: 0}).Return(nil, errors.New("db error"))

	_, err := svc.List(context.Background(), dto.ListUserSessionsRequest{UserID: uid, Limit: 1, Offset: 0})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserSessionService_ListActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	sessions := []db.UserSession{{Status: "active"}}
	mockRepo.EXPECT().ListActiveSessions(gomock.Any(), gomock.Any()).Return(sessions, nil)

	got, err := svc.ListActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(sessions) {
		t.Fatalf("got %d want %d", len(got), len(sessions))
	}
}

func TestUserSessionService_ListActive_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	mockRepo.EXPECT().ListActiveSessions(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

	_, err := svc.ListActive(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserSessionService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var id pgtype.UUID
	var endedAt pgtype.Timestamptz
	up := db.UpdateUserSessionParams{ID: id, Status: "active", EndedAt: endedAt}
	want := db.UserSession{ID: id, Status: "active"}

	mockRepo.EXPECT().UpdateUserSession(gomock.Any(), up).Return(want, nil)

	got, err := svc.Update(context.Background(), dto.UpdateUserSessionRequest{ID: id, Status: up.Status, EndedAt: endedAt})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != want.Status {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserSessionService_Update_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var id pgtype.UUID
	var endedAt pgtype.Timestamptz
	up := db.UpdateUserSessionParams{ID: id, Status: "active", EndedAt: endedAt}

	mockRepo.EXPECT().UpdateUserSession(gomock.Any(), up).Return(db.UserSession{}, errors.New("db error"))

	_, err := svc.Update(context.Background(), dto.UpdateUserSessionRequest{ID: id, Status: up.Status, EndedAt: endedAt})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserSessionService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().DeleteUserSession(gomock.Any(), id).Return(nil)

	err := svc.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserSessionService_Delete_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserSessionRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserSessionService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().DeleteUserSession(gomock.Any(), id).Return(errors.New("db error"))

	err := svc.Delete(context.Background(), id)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
