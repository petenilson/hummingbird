package postgres_test

import (
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
	db, err := MustOpenTestDB(m)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	db.Now = func() time.Time {
		return time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
	}

	exit_code := m.Run()

	db.Close()
	os.Exit(exit_code)
}

func NewTestDB(tb testing.TB) *postgres.DB {
	tb.Helper()
	container := MustCreateContainer(tb)
	dsn, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		tb.Fatal(err)
	}
	return postgres.NewDB(dsn)
}

func MustCreateContainer(tb testing.TB) *container_postgres.PostgresContainer {
	tb.Helper()
	container, err := container_postgres.Run(
		context.Background(),
		"docker.io/postgres:16-alpine",
		container_postgres.WithDatabase("TestDB"),
		container_postgres.WithUsername("TestUser"),
		container_postgres.WithPassword("TestPassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		tb.Fatal(err)
	}
	return container
}

func MustOpenTestDB(tb testing.TB) *postgres.DB {
	tb.Helper()
	ctx := context.Background()

	container := MustCreateContainer(tb)
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		tb.Fatal(err)
	}
	db := postgres.NewDB(connStr)

	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}
	return db
}
