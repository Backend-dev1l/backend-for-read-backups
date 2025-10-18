package db

import (
	"context"
	"testing"

	"test-http/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	t.Run("Create/Get/Update/Delete user", func(t *testing.T) {
		userArg := CreateUserParams{
			Username: "newuser",
			Email:    "newuser@example.com",
		}
		created, err := queries.CreateUser(ctx, userArg)
		require.NoError(t, err)
		require.Equal(t, userArg.Username, created.Username)
		require.Equal(t, userArg.Email, created.Email)

		retrieved, err := queries.GetUser(ctx, created.ID)
		require.NoError(t, err)
		require.Equal(t, created.ID, retrieved.ID)

		updateArg := UpdateUserParams{
			ID:       created.ID,
			Username: "updateduser",
			Email:    "updated@example.com",
		}
		updated, err := queries.UpdateUser(ctx, updateArg)
		require.NoError(t, err)
		require.Equal(t, updateArg.Username, updated.Username)
		require.Equal(t, updateArg.Email, updated.Email)

		err = queries.DeleteUser(ctx, created.ID)
		require.NoError(t, err)
		_, err = queries.GetUser(ctx, created.ID)
		require.Error(t, err)
	})

	t.Run("UNIQUE violation", func(t *testing.T) {
		user := CreateUserParams{Username: "user1", Email: "test1@example.com"}
		_, err := queries.CreateUser(ctx, user)
		require.Error(t, err)
	})

	t.Run("NULL violation", func(t *testing.T) {
		user := CreateUserParams{Username: "", Email: ""}
		_, err := queries.CreateUser(ctx, user)
		require.Error(t, err)
	})

	// Проверка на race condition: параллельное создание с одинаковым email
	t.Run("Race: duplicate email", func(t *testing.T) {
		email := "race@example.com"
		done := make(chan error, 2)
		go func() {
			_, err := queries.CreateUser(ctx, CreateUserParams{Username: uuid.New().String(), Email: email})
			done <- err
		}()
		go func() {
			_, err := queries.CreateUser(ctx, CreateUserParams{Username: uuid.New().String(), Email: email})
			done <- err
		}()
		err1 := <-done
		err2 := <-done
		require.True(t, err1 != nil || err2 != nil, "Expected UNIQUE constraint violation in one of the goroutines")
	})
}

func TestUserRepository_CRUDEdgeCases(t *testing.T) {
	t.Parallel()
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)
	testutil.CleanupDB(t, pool, "users")
	testutil.SeedTestData(t, pool, "../../testdata/fixtures")
	queries := New(pool)
	ctx := context.Background()

	t.Run("UNIQUE violation", func(t *testing.T) {
		user := CreateUserParams{Username: "Test User 1", Email: "test1@example.com"}
		_, err := queries.CreateUser(ctx, user)
		require.Error(t, err, "Should return UNIQUE error")
	})
	t.Run("NULL violation", func(t *testing.T) {
		user := CreateUserParams{Username: "", Email: ""}
		_, err := queries.CreateUser(ctx, user)
		require.Error(t, err, "Should return NOT NULL error")
	})
	t.Run("Race condition: duplicate email", func(t *testing.T) {
		email := uuid.New().String() + "@race.com"
		done := make(chan error, 2)
		go func() {
			_, err := queries.CreateUser(ctx, CreateUserParams{Username: "u1", Email: email})
			done <- err
		}()
		go func() {
			_, err := queries.CreateUser(ctx, CreateUserParams{Username: "u2", Email: email})
			done <- err
		}()
		err1 := <-done
		err2 := <-done
		require.True(t, err1 != nil || err2 != nil, "Expected UNIQUE constraint violation in race condition")
	})
	// Users has no FK or CHECK constraints to test
}
