package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	db "test-http/internal/db"
	"test-http/internal/dto"
	mocks "test-http/internal/db/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// testHelper contains common test utilities
type statisticsTestHelper struct {
	ctrl       *gomock.Controller
	mockRepo   *mocks.MockUserStatisticsRepo
	service    *UserStatisticsService
	logger     *slog.Logger
	ctx        context.Context
}

// newStatisticsTestHelper creates a new test helper with initialized mocks
func newStatisticsTestHelper(t *testing.T) *statisticsTestHelper {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockUserStatisticsRepo(ctrl)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	service := NewUserStatisticsService(mockRepo, logger)

	return &statisticsTestHelper{
		ctrl:     ctrl,
		mockRepo: mockRepo,
		service:  service,
		logger:   logger,
		ctx:      context.Background(),
	}
}

// cleanup performs cleanup after test
func (h *statisticsTestHelper) cleanup() {
	h.ctrl.Finish()
}

// newStatisticsUUID creates a new pgtype.UUID from string
func newStatisticsUUID(t *testing.T, uuidStr string) pgtype.UUID {
	t.Helper()
	
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(uuidStr)
	if err != nil {
		t.Fatalf("failed to scan UUID: %v", err)
	}
	
	return pgUUID
}

// randomStatisticsUUID generates a random pgtype.UUID
func randomStatisticsUUID(t *testing.T) pgtype.UUID {
	t.Helper()
	return newStatisticsUUID(t, uuid.New().String())
}

// newNumeric creates a pgtype.Numeric from float64
func newNumeric(t *testing.T, value float64) pgtype.Numeric {
	t.Helper()
	
	var num pgtype.Numeric
	err := num.Scan(fmt.Sprintf("%f", value))
	if err != nil {
		t.Fatalf("failed to create numeric: %v", err)
	}
	
	return num
}

func TestUserStatisticsService_Create(t *testing.T) {
	t.Run("successfully creates user statistics", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		accuracy := newNumeric(t, 85.5)
		
		request := dto.CreateStatisticsRequest{
			UserID:            userID,
			TotalWordsLearned: 100,
			Accuracy:          accuracy,
			TotalTime:         3600,
		}

		expectedParams := db.CreateUserStatisticsParams{
			UserID:            userID,
			TotalWordsLearned: 100,
			Accuracy:          accuracy,
			TotalTime:         3600,
		}

		expectedStats := db.UserStatistic{
			UserID:            userID,
			TotalWordsLearned: 100,
			Accuracy:          accuracy,
			TotalTime:         3600,
		}

		h.mockRepo.EXPECT().
			CreateUserStatistics(h.ctx, expectedParams).
			Return(expectedStats, nil)

		result, err := h.service.Create(h.ctx, request)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if result.UserID.Bytes != expectedStats.UserID.Bytes {
			t.Errorf("expected UserID %v, got %v", expectedStats.UserID, result.UserID)
		}

		if result.TotalWordsLearned != expectedStats.TotalWordsLearned {
			t.Errorf("expected TotalWordsLearned %d, got %d", expectedStats.TotalWordsLearned, result.TotalWordsLearned)
		}

		expectedAccuracyFloat, _ := expectedStats.Accuracy.Float64Value()
		resultAccuracyFloat, _ := result.Accuracy.Float64Value()
		if resultAccuracyFloat.Float64 != expectedAccuracyFloat.Float64 {
			t.Errorf("expected Accuracy %v, got %v", expectedAccuracyFloat.Float64, resultAccuracyFloat.Float64)
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		request := dto.CreateStatisticsRequest{
			UserID:            userID,
			TotalWordsLearned: 50,
		}

		expectedError := errors.New("database connection failed")

		h.mockRepo.EXPECT().
			CreateUserStatistics(h.ctx, gomock.Any()).
			Return(db.UserStatistic{}, expectedError)

		_, err := h.service.Create(h.ctx, request)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserStatisticsService_GetByID(t *testing.T) {
	t.Run("successfully retrieves user statistics by user ID", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		request := dto.GetStatisticsRequest{UserID: userID}
		accuracy := newNumeric(t, 90.0)

		expectedStats := db.UserStatistic{
			UserID:            userID,
			TotalWordsLearned: 75,
			Accuracy:          accuracy,
			TotalTime:         2400,
		}

		h.mockRepo.EXPECT().
			GetUserStatistics(h.ctx, userID).
			Return(expectedStats, nil)

		result, err := h.service.GetByID(h.ctx, request)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if result.UserID.Bytes != expectedStats.UserID.Bytes {
			t.Errorf("expected UserID %v, got %v", expectedStats.UserID, result.UserID)
		}

		if result.TotalWordsLearned != expectedStats.TotalWordsLearned {
			t.Errorf("expected TotalWordsLearned %d, got %d", expectedStats.TotalWordsLearned, result.TotalWordsLearned)
		}
	})

	t.Run("returns error when user statistics not found", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		request := dto.GetStatisticsRequest{UserID: userID}

		expectedError := errors.New("statistics not found")

		h.mockRepo.EXPECT().
			GetUserStatistics(h.ctx, userID).
			Return(db.UserStatistic{}, expectedError)

		_, err := h.service.GetByID(h.ctx, request)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserStatisticsService_List(t *testing.T) {
	t.Run("successfully lists all user statistics", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		request := dto.ListStatisticsRequest{
			Limit:  10,
			Offset: 0,
		}

		expectedParams := db.ListUserStatisticsParams{
			Limit:  10,
			Offset: 0,
		}

		userID1 := randomStatisticsUUID(t)
		userID2 := randomStatisticsUUID(t)
		accuracy1 := newNumeric(t, 80.0)
		accuracy2 := newNumeric(t, 95.0)

		expectedStats := []db.UserStatistic{
			{UserID: userID1, TotalWordsLearned: 50, Accuracy: accuracy1},
			{UserID: userID2, TotalWordsLearned: 100, Accuracy: accuracy2},
		}

		h.mockRepo.EXPECT().
			ListUserStatistics(h.ctx, expectedParams).
			Return(expectedStats, nil)

		result, err := h.service.List(h.ctx, request)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if len(result) != len(expectedStats) {
			t.Errorf("expected %d statistics, got %d", len(expectedStats), len(result))
		}
	})

	t.Run("returns error when repository fails to list", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		request := dto.ListStatisticsRequest{
			Limit:  10,
			Offset: 0,
		}

		expectedError := errors.New("failed to query database")

		h.mockRepo.EXPECT().
			ListUserStatistics(h.ctx, gomock.Any()).
			Return(nil, expectedError)

		_, err := h.service.List(h.ctx, request)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns empty list when no statistics exist", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		request := dto.ListStatisticsRequest{
			Limit:  10,
			Offset: 0,
		}

		h.mockRepo.EXPECT().
			ListUserStatistics(h.ctx, gomock.Any()).
			Return([]db.UserStatistic{}, nil)

		result, err := h.service.List(h.ctx, request)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty list, got %d items", len(result))
		}
	})
}

func TestUserStatisticsService_Update(t *testing.T) {
	t.Run("successfully updates user statistics", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		accuracy := newNumeric(t, 92.5)
		
		request := dto.UpdateStatisticsRequest{
			UserID:            userID,
			TotalWordsLearned: 150,
			Accuracy:          accuracy,
			TotalTime:         5000,
		}

		expectedParams := db.UpdateUserStatisticsParams{
			UserID:            userID,
			TotalWordsLearned: 150,
			Accuracy:          accuracy,
			TotalTime:         5000,
		}

		expectedStats := db.UserStatistic{
			UserID:            userID,
			TotalWordsLearned: 150,
			Accuracy:          accuracy,
			TotalTime:         5000,
		}

		h.mockRepo.EXPECT().
			UpdateUserStatistics(h.ctx, expectedParams).
			Return(expectedStats, nil)

		result, err := h.service.Update(h.ctx, request)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if result.TotalWordsLearned != expectedStats.TotalWordsLearned {
			t.Errorf("expected TotalWordsLearned %d, got %d", expectedStats.TotalWordsLearned, result.TotalWordsLearned)
		}

		expectedAccuracyFloat, _ := expectedStats.Accuracy.Float64Value()
		resultAccuracyFloat, _ := result.Accuracy.Float64Value()
		if resultAccuracyFloat.Float64 != expectedAccuracyFloat.Float64 {
			t.Errorf("expected Accuracy %v, got %v", expectedAccuracyFloat.Float64, resultAccuracyFloat.Float64)
		}
	})

	t.Run("returns error when repository update fails", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		request := dto.UpdateStatisticsRequest{
			UserID:            userID,
			TotalWordsLearned: 200,
		}

		expectedError := errors.New("update operation failed")

		h.mockRepo.EXPECT().
			UpdateUserStatistics(h.ctx, gomock.Any()).
			Return(db.UserStatistic{}, expectedError)

		_, err := h.service.Update(h.ctx, request)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserStatisticsService_Delete(t *testing.T) {
	t.Run("successfully deletes user statistics", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		request := dto.DeleteStatisticsRequest{UserID: userID}

		h.mockRepo.EXPECT().
			DeleteUserStatistics(h.ctx, userID).
			Return(nil)

		err := h.service.Delete(h.ctx, request)

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("returns error when repository delete fails", func(t *testing.T) {
		h := newStatisticsTestHelper(t)
		defer h.cleanup()

		userID := randomStatisticsUUID(t)
		request := dto.DeleteStatisticsRequest{UserID: userID}

		expectedError := errors.New("delete operation failed")

		h.mockRepo.EXPECT().
			DeleteUserStatistics(h.ctx, userID).
			Return(expectedError)

		err := h.service.Delete(h.ctx, request)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
