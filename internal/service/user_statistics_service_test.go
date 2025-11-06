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

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserStatisticsService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	params := CreateUserStatisticsParams{UserID: uid, TotalWordsLearned: 10}
	want := db.UserStatistic{UserID: uid, TotalWordsLearned: params.TotalWordsLearned}

	mockRepo.EXPECT().CreateUserStatistics(gomock.Any(), db.CreateUserStatisticsParams{
		UserID:            uid,
		TotalWordsLearned: params.TotalWordsLearned,
	}).Return(want, nil)

	got, err := svc.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserStatisticsService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	params := CreateUserStatisticsParams{UserID: uid}

	mockRepo.EXPECT().CreateUserStatistics(gomock.Any(), gomock.Any()).Return(db.UserStatistic{}, errors.New("fail"))

	_, err := svc.Create(context.Background(), params)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserStatisticsService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	want := db.UserStatistic{UserID: uid}

	mockRepo.EXPECT().GetUserStatistics(gomock.Any(), uid).Return(want, nil)

	got, err := svc.GetByID(context.Background(), uid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserStatisticsService_GetByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	mockRepo.EXPECT().GetUserStatistics(gomock.Any(), uid).Return(db.UserStatistic{}, errors.New("fail"))

	_, err := svc.GetByID(context.Background(), uid)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserStatisticsService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	params := UpdateUserStatisticsParams{UserID: uid, TotalWordsLearned: 20}
	want := db.UserStatistic{UserID: uid, TotalWordsLearned: params.TotalWordsLearned}

	mockRepo.EXPECT().UpdateUserStatistics(gomock.Any(), db.UpdateUserStatisticsParams{
		UserID:            uid,
		TotalWordsLearned: params.TotalWordsLearned,
	}).Return(want, nil)

	got, err := svc.Update(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUserStatisticsService_Update_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	mockRepo.EXPECT().UpdateUserStatistics(gomock.Any(), gomock.Any()).Return(db.UserStatistic{}, errors.New("fail"))

	_, err := svc.Update(context.Background(), UpdateUserStatisticsParams{UserID: uid})
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}

func TestUserStatisticsService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	mockRepo.EXPECT().DeleteUserStatistics(gomock.Any(), uid).Return(nil)

	err := svc.Delete(context.Background(), uid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserStatisticsService_Delete_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	svc := NewUserStatisticsService(mockRepo, logger)

	var uid pgtype.UUID
	mockRepo.EXPECT().DeleteUserStatistics(gomock.Any(), uid).Return(errors.New("fail"))

	err := svc.Delete(context.Background(), uid)
	if err == nil {
		t.Fatalf("expected error from repo, got nil")
	}
}
