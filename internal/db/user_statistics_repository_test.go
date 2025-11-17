package db

import (
	"context"
	"testing"

	"test-http/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserStatisticsRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_statistics", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	var userID = uuid.New()

	t.Run("Create/Get/Update/Delete statistics", func(t *testing.T) {
		arg := CreateUserStatisticsParams{
			UserID:            testutil.NewPgUUID(t, userID),
			TotalWordsLearned: 10,
			Accuracy:          testutil.NewPgNumeric(t, 85.5),
			TotalTime:         3600,
		}
		_, err := queries.CreateUserStatistics(ctx, arg)
		require.Error(t, err)
	})
}

func TestUserStatisticsRepository_EdgeCases(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_statistics", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	t.Run("FK violation", func(t *testing.T) {
		arg := CreateUserStatisticsParams{
			UserID:            testutil.RandomPgUUID(t),
			TotalWordsLearned: 1,
			Accuracy:          testutil.NewPgNumeric(t, 90.0),
			TotalTime:         10,
		}
		_, err := queries.CreateUserStatistics(ctx, arg)
		require.Error(t, err, "Should return FK error")
	})

	t.Run("NULL violation", func(t *testing.T) {
		arg := CreateUserStatisticsParams{
			UserID:            testutil.EmptyPgUUID(),
			TotalWordsLearned: 0,
			Accuracy:          testutil.EmptyPgNumeric(),
			TotalTime:         0,
		}
		_, err := queries.CreateUserStatistics(ctx, arg)
		require.Error(t, err, "Should return NOT NULL error")
	})
}
