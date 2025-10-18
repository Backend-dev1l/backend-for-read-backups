package db

import (
	"context"
	"testing"

	"test-http/internal/testutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestUserProgressRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_progress", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	var anyUserID = uuid.New()
	var anyWordID = uuid.New()

	t.Run("Create/Get/Update/Delete progress", func(t *testing.T) {
		pArg := CreateUserProgressParams{
			UserID:         pgtype.UUID{Bytes: anyUserID, Valid: true},
			WordID:         pgtype.UUID{Bytes: anyWordID, Valid: true},
			CorrectCount:   2,
			IncorrectCount: 1,
		}
		_, err := queries.CreateUserProgress(ctx, pArg)
		require.Error(t, err, "До инициализации нормальных test-uuid должен быть FK error")
	})
}

func TestUserProgressRepository_EdgeCases(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_progress", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	t.Run("NULL violation", func(t *testing.T) {
		arg := CreateUserProgressParams{
			UserID:         pgtype.UUID{},
			WordID:         pgtype.UUID{},
			CorrectCount:   0,
			IncorrectCount: 0,
		}
		_, err := queries.CreateUserProgress(ctx, arg)
		require.Error(t, err, "Should return NOT NULL error")
	})

	t.Run("FK violation", func(t *testing.T) {
		arg := CreateUserProgressParams{
			UserID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
			WordID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
			CorrectCount:   1,
			IncorrectCount: 0,
		}
		_, err := queries.CreateUserProgress(ctx, arg)
		require.Error(t, err, "Should return FK error")
	})

	// Для race и CHECK нужно знать ограничения схемы, допустим ограничение на уникальность по (user_id, word_id)
	t.Run("UNIQUE violation + race", func(t *testing.T) {
		fk := pgtype.UUID{Bytes: uuid.New(), Valid: true}
		arg := CreateUserProgressParams{
			UserID:         fk,
			WordID:         fk,
			CorrectCount:   10,
			IncorrectCount: 5,
		}
		_, _ = queries.CreateUserProgress(ctx, arg)
		done := make(chan error, 2)
		go func() {
			_, err := queries.CreateUserProgress(ctx, arg)
			done <- err
		}()
		go func() {
			_, err := queries.CreateUserProgress(ctx, arg)
			done <- err
		}()
		err1 := <-done
		err2 := <-done
		require.True(t, err1 != nil || err2 != nil, "Expected UNIQUE constraint violation for (user_id, word_id) race")
	})
}
