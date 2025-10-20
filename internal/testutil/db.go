package testutil

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx/v5/pgxpool"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pgContainer *postgres.PostgresContainer

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:17"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("app_user"),
		postgres.WithPassword("my-password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
		testcontainers.WithExposedPorts("5432/tcp"),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	pgContainer = container

	uri, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to get connection string: %v", err)
	}

	applyMigrations(t, uri)

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return pool
}

func TeardownTestDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	if pool != nil {
		pool.Close()
	}
	if pgContainer != nil {
		_ = pgContainer.Terminate(context.Background())
		pgContainer = nil
	}
}

func SeedTestData(t *testing.T, pool *pgxpool.Pool, fixturesPath string) {
	t.Helper()
	ctx := context.Background()

	// Users
	usersPath := fmt.Sprintf("%s/users.json", fixturesPath)
	userData, err := ioutil.ReadFile(usersPath)
	if err != nil {
		t.Fatalf("failed to read users fixture: %v", err)
	}
	var users []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err = json.Unmarshal(userData, &users); err != nil {
		t.Fatalf("failed to unmarshal users fixture: %v", err)
	}
	for _, u := range users {
		_, err := pool.Exec(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", u.ID, u.Name, u.Email)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}
	}

	progPath := fmt.Sprintf("%s/user_progress.json", fixturesPath)
	progData, err := ioutil.ReadFile(progPath)
	if err != nil {
		t.Fatalf("failed to read user_progress fixture: %v", err)
	}
	var progress []struct {
		ID       int `json:"id"`
		UserID   int `json:"user_id"`
		Progress int `json:"progress"`
	}
	if err = json.Unmarshal(progData, &progress); err != nil {
		t.Fatalf("failed to unmarshal user_progress fixture: %v", err)
	}
	for _, p := range progress {
		_, err := pool.Exec(ctx, "INSERT INTO user_progress (id, user_id, progress) VALUES ($1, $2, $3)", p.ID, p.UserID, p.Progress)
		if err != nil {
			t.Fatalf("failed to insert user_progress: %v", err)
		}
	}
}

func applyMigrations(t *testing.T, dsn string) {
	t.Helper()

	// Сначала перегенерируем хэши для основной директории миграций
	hashCmd := exec.Command("atlas", "migrate", "hash", "--dir", "file://migrations")
	hashCmd.Stdout = os.Stdout
	hashCmd.Stderr = os.Stderr
	if err := hashCmd.Run(); err != nil {
		t.Fatalf("failed to hash migrations: %v", err)
	}

	// Затем применяем миграции с указанием правильного пути
	cmd := exec.Command("atlas", "migrate", "apply",
		"--dir", "file://migrations",
		"--url", dsn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}
}

func CleanupDB(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()
	ctx := context.Background()
	if len(tables) == 0 {
		t.Fatal("No tables provided for cleanup")
	}
	q := "TRUNCATE TABLE "
	for i, tbl := range tables {
		q += tbl
		if i < len(tables)-1 {
			q += ", "
		}
	}
	q += " CASCADE;"
	_, err := pool.Exec(ctx, q)
	if err != nil {
		t.Fatalf("failed to cleanup db: %v", err)
	}
}
