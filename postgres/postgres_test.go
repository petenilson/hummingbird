package postgres_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"context"

	"github.com/petenilson/go-ledger/postgres"

	"github.com/testcontainers/testcontainers-go"
	container_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var DB *postgres.DB

func TestMain(m *testing.M) {
	// Setup code...
	db, err := MustReturnTestDB()
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	db.Now = func() time.Time {
		return time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
	}

	exit_code := m.Run()

	// Teardown code...
	db.Close()
	os.Exit(exit_code)
}

// TestDB allows us to define at test time what will be returned
// when Begin is called on the DB. Our goal is to have begin return
// an inner pseudo transaction that can be reverted after every
// test function.
type TestDB struct {
	db      *postgres.DB
	BeginFn func(context.Context) (*postgres.Tx, error)
	outerTx *postgres.Tx
	Now     func() time.Time
}

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

func MustReturnTestDB() (*postgres.DB, error) {
	ctx := context.Background()

	container, err := MustStartContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed Creating Test Container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Failed Getting Connection String: %w", err)
	}
	db := postgres.NewDB(connStr)

	if err := db.Open(); err != nil {
		return nil, fmt.Errorf("Failed To Open DB %s", err)
	}

	return db, nil
}
