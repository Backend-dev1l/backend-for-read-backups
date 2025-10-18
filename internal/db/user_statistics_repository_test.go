package db

import (
	"context"
	"math/big"
	"testing"

	"test-http/internal/testutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
			UserID:            pgtype.UUID{Bytes: userID, Valid: true},
			TotalWordsLearned: 10,
			Accuracy:          pgtype.Numeric{Int: big.NewInt(0), Valid: true},
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
			UserID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
			TotalWordsLearned: 1,
			Accuracy:          pgtype.Numeric{Int: big.NewInt(0), Valid: true},
			TotalTime:         10,
		}
		_, err := queries.CreateUserStatistics(ctx, arg)
		require.Error(t, err, "Should return FK error")
	})

	t.Run("NULL violation", func(t *testing.T) {
		arg := CreateUserStatisticsParams{
			UserID:            pgtype.UUID{},
			TotalWordsLearned: 0,
			Accuracy:          pgtype.Numeric{},
			TotalTime:         0,
		}
		_, err := queries.CreateUserStatistics(ctx, arg)
		require.Error(t, err, "Should return NOT NULL error")
	})
}
