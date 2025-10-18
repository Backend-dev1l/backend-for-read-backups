package db

import (
	"context"
	"testing"

	"test-http/internal/testutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestUserWordSetsRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_word_sets", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	var userID = uuid.New()
	var setID = uuid.New()

	t.Run("Create/Get/Delete word set", func(t *testing.T) {
		arg := CreateUserWordSetParams{
			UserID:    pgtype.UUID{Bytes: userID, Valid: true},
			WordSetID: pgtype.UUID{Bytes: setID, Valid: true},
		}
		_, err := queries.CreateUserWordSet(ctx, arg)
		require.Error(t, err)
	})
}

func TestUserWordSetsRepository_EdgeCases(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_word_sets", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	t.Run("FK violation", func(t *testing.T) {
		arg := CreateUserWordSetParams{
			UserID:    pgtype.UUID{Bytes: uuid.New(), Valid: true},
			WordSetID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		}
		_, err := queries.CreateUserWordSet(ctx, arg)
		require.Error(t, err, "Should return FK error")
	})

	t.Run("NULL violation", func(t *testing.T) {
		arg := CreateUserWordSetParams{
			UserID:    pgtype.UUID{},
			WordSetID: pgtype.UUID{},
		}
		_, err := queries.CreateUserWordSet(ctx, arg)
		require.Error(t, err, "Should return NOT NULL error")
	})
}
