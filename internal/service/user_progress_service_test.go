package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	db "test-http/internal/db"
	"test-http/internal/dto"
	mockdb "test-http/internal/db/mocks"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserProgressService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var uid, wid pgtype.UUID
	params := dto.CreateUserProgressRequest{UserID: uid, WordID: wid, CorrectCount: 1, IncorrectCount: 0}
	want := db.UserProgress{UserID: uid, WordID: wid, CorrectCount: params.CorrectCount, IncorrectCount: params.IncorrectCount}

	mockRepo.EXPECT().CreateUserProgress(gomock.Any(), db.CreateUserProgressParams{
		UserID:         uid,
		WordID:         wid,
		CorrectCount:   params.CorrectCount,
		IncorrectCount: params.IncorrectCount,
	}).Return(want, nil)

	got, err := svc.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != want.UserID || got.WordID != want.WordID {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserProgressService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var uid, wid pgtype.UUID
	params := dto.CreateUserProgressRequest{UserID: uid, WordID: wid, CorrectCount: 1, IncorrectCount: 0}

	mockRepo.EXPECT().CreateUserProgress(gomock.Any(), gomock.Any()).Return(db.UserProgress{}, errors.New("fail"))

	_, err := svc.Create(context.Background(), params)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserProgressService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var id pgtype.UUID
	want := db.UserProgress{ID: id}
	mockRepo.EXPECT().GetUserProgress(gomock.Any(), id).Return(want, nil)

	got, err := svc.GetByID(context.Background(), dto.GetUserProgressRequest{ID: id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserProgressService_GetByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().GetUserProgress(gomock.Any(), id).Return(db.UserProgress{}, errors.New("not found"))

	_, err := svc.GetByID(context.Background(), dto.GetUserProgressRequest{ID: id})
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserProgressService_GetByUserAndWord_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var uid, wid pgtype.UUID
	params := db.GetUserProgressByUserAndWordParams{UserID: uid, WordID: wid}
	want := db.UserProgress{UserID: uid, WordID: wid}
	mockRepo.EXPECT().GetUserProgressByUserAndWord(gomock.Any(), params).Return(want, nil)

	got, err := svc.GetByUserAndWord(context.Background(), dto.GetUserProgressByUserAndWordRequest{UserID: uid, WordID: wid})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != want.UserID || got.WordID != want.WordID {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserProgressService_GetByUserAndWord_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var uid, wid pgtype.UUID
	mockRepo.EXPECT().GetUserProgressByUserAndWord(gomock.Any(), gomock.Any()).Return(db.UserProgress{}, errors.New("fail"))

	_, err := svc.GetByUserAndWord(context.Background(), dto.GetUserProgressByUserAndWordRequest{UserID: uid, WordID: wid})
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserProgressService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var uid pgtype.UUID
	filters := dto.ListUserProgressRequest{UserID: uid, Limit: 10, Offset: 0}
	want := []db.UserProgress{{UserID: uid}}
	mockRepo.EXPECT().ListUserProgress(gomock.Any(), uid).Return(want, nil)

	got, err := svc.List(context.Background(), filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d want %d", len(got), len(want))
	}
}

func TestUserProgressService_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var uid pgtype.UUID
	filters := dto.ListUserProgressRequest{UserID: uid}
	mockRepo.EXPECT().ListUserProgress(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))

	_, err := svc.List(context.Background(), filters)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserProgressService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var id pgtype.UUID
	params := dto.UpdateUserProgressRequest{ID: id, CorrectCount: 2, IncorrectCount: 1}
	want := db.UserProgress{ID: id, CorrectCount: params.CorrectCount, IncorrectCount: params.IncorrectCount}

	mockRepo.EXPECT().UpdateUserProgress(gomock.Any(), db.UpdateUserProgressParams{ID: id, CorrectCount: params.CorrectCount, IncorrectCount: params.IncorrectCount}).Return(want, nil)

	got, err := svc.Update(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserProgressService_Update_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var id pgtype.UUID
	params := dto.UpdateUserProgressRequest{ID: id}
	mockRepo.EXPECT().UpdateUserProgress(gomock.Any(), gomock.Any()).Return(db.UserProgress{}, errors.New("fail"))

	_, err := svc.Update(context.Background(), params)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserProgressService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().DeleteUserProgress(gomock.Any(), id).Return(nil)

	err := svc.Delete(context.Background(), dto.DeleteUserProgressRequest{ID: id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserProgressService_Delete_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockdb.NewMockUserProgressRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserProgressService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().DeleteUserProgress(gomock.Any(), id).Return(errors.New("fail"))

	err := svc.Delete(context.Background(), dto.DeleteUserProgressRequest{ID: id})
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}
