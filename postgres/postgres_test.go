package postgres_test

import (
	"testing"
	"time"

	"context"

	"github.com/petenilson/go-ledger/postgres"

	"github.com/testcontainers/testcontainers-go"
	container_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func MustStartContainer(ctx context.Context) (*container_postgres.PostgresContainer, error) {
	if postgresContainer, err := container_postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		container_postgres.WithDatabase("TestDB"),
		container_postgres.WithUsername("TestUser"),
		container_postgres.WithPassword("TestPassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	); err != nil {
		return nil, err
	} else {
		return postgresContainer, nil
	}
}

func MustOpenDB(t *testing.T) (*postgres.DB, func()) {
	ctx := context.Background()

	container, err := MustStartContainer(ctx)
	if err != nil {
		t.Fatalf("Failed Creating Test Container")
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed Getting Connection String")
	}
	db := postgres.NewDB(connStr)
	db.Now = func() time.Time {
		return time.Date(2024, time.February, 1, 12, 01, 03, 0, time.UTC)
	}

	if err := db.Open(); err != nil {
		t.Fatalf("Failed To Open DB %s", err)
	}

	return db, func() {
		container.Terminate(ctx)
		db.Close()
	}
}
