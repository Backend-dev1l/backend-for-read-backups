package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"testing"

	db "test-http/internal/db"
	mocks "test-http/internal/db/mocks"
	"test-http/internal/dto"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserWordSetService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserWordSetRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserWordSetService(mockRepo, logger)

	var uid, wid pgtype.UUID
	want := db.UserWordSet{UserID: uid, WordSetID: wid}

	mockRepo.EXPECT().CreateUserWordSet(gomock.Any(), db.CreateUserWordSetParams{
		UserID:    uid,
		WordSetID: wid,
	}).Return(want, nil)

	got, err := svc.Create(context.Background(), dto.CreateUserWordSetRequest{UserID: uid, WordSetID: wid})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserWordSetService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserWordSetRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserWordSetService(mockRepo, logger)

	var uid, wid pgtype.UUID
	mockRepo.EXPECT().CreateUserWordSet(gomock.Any(), gomock.Any()).Return(db.UserWordSet{}, errors.New("fail"))

	_, err := svc.Create(context.Background(), dto.CreateUserWordSetRequest{UserID: uid, WordSetID: wid})
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserWordSetService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserWordSetRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserWordSetService(mockRepo, logger)

	var id pgtype.UUID
	want := db.UserWordSet{ID: id}

	mockRepo.EXPECT().GetUserWordSet(gomock.Any(), id).Return(want, nil)

	got, err := svc.GetByID(context.Background(), dto.GetUserWordSetRequest{ID: id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserWordSetService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserWordSetRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserWordSetService(mockRepo, logger)

	var uid pgtype.UUID
	want := []db.UserWordSet{{UserID: uid}}
	mockRepo.EXPECT().ListUserWordSets(gomock.Any(), db.ListUserWordSetsParams{
		UserID: uid,
		Limit:  10,
		Offset: 0,
	}).Return(want, nil)

	got, err := svc.List(context.Background(), dto.ListUserWordSetsRequest{UserID: uid, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d want %d", len(got), len(want))
	}
}

func TestUserWordSetService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserWordSetRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserWordSetService(mockRepo, logger)

	var id pgtype.UUID
	mockRepo.EXPECT().DeleteUserWordSet(gomock.Any(), id).Return(nil)

	err := svc.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
