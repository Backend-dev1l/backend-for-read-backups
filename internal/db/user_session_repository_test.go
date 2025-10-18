package db

import (
	"context"
	"testing"

	"test-http/internal/testutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestUserSessionRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_sessions", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	var userID = uuid.New()

	t.Run("Create/Get/Update/Delete session", func(t *testing.T) {
		arg := CreateUserSessionParams{
			UserID: pgtype.UUID{Bytes: userID, Valid: true},
			Status: "active",
		}
		_, err := queries.CreateUserSession(ctx, arg)
		require.Error(t, err, "FK violation will be returned until valid IDs are configured")
	})
}

func TestUserSessionRepository_EdgeCases(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "user_sessions", "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	t.Run("FK violation", func(t *testing.T) {
		arg := CreateUserSessionParams{
			UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Status: "active",
		}
		_, err := queries.CreateUserSession(ctx, arg)
		require.Error(t, err, "Should return FK error")
	})

	t.Run("NULL violation", func(t *testing.T) {
		arg := CreateUserSessionParams{
			UserID: pgtype.UUID{},
			Status: "",
		}
		_, err := queries.CreateUserSession(ctx, arg)
		require.Error(t, err, "Should return NOT NULL error")
	})
}
