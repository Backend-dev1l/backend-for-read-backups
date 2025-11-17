package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pgContainer testcontainers.Container

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	// Create and start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image: "postgres:17",
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "my-password",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	// Привязываем к глобальной переменной только интерфейс Container
	pgContainer = container

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	user := "testuser"
	password := "my-password"
	database := "testdb"
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port.Port(), database)

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
	userData, err := os.ReadFile(usersPath)
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
	progData, err := os.ReadFile(progPath)
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

// NewPgUUID creates a pgtype.UUID from uuid.UUID
func NewPgUUID(t *testing.T, id uuid.UUID) pgtype.UUID {
	t.Helper()
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(id.String()); err != nil {
		t.Fatalf("failed to create pgtype.UUID: %v", err)
	}
	return pgUUID
}

// RandomPgUUID creates a random pgtype.UUID
func RandomPgUUID(t *testing.T) pgtype.UUID {
	t.Helper()
	return NewPgUUID(t, uuid.New())
}

// NewPgNumeric creates a pgtype.Numeric from float64
func NewPgNumeric(t *testing.T, value float64) pgtype.Numeric {
	t.Helper()
	var num pgtype.Numeric
	if err := num.Scan(value); err != nil {
		t.Fatalf("failed to create pgtype.Numeric: %v", err)
	}
	return num
}

// EmptyPgUUID creates an invalid/NULL pgtype.UUID
func EmptyPgUUID() pgtype.UUID {
	return pgtype.UUID{}
}

// EmptyPgNumeric creates an invalid/NULL pgtype.Numeric
func EmptyPgNumeric() pgtype.Numeric {
	return pgtype.Numeric{}
}
