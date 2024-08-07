package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/cmd/hummingbird"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func MustRunMain(tb testing.TB) *main.Main {
	tb.Helper()

	container := MustCreateContainer(tb)
	dsn, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		tb.Fatal(err)
	}

	m := main.NewMain()
	m.Config.DB.DSN = dsn
	m.Config.HTTP.Address = "localhost:8000"

	if err := m.Run(context.Background()); err != nil {
		tb.Fatal(err)
	}

	m.DB.Now = func() time.Time { return time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC) }
	return m
}

func MustCreateAccount(tb testing.TB, m *main.Main, account *hummingbird.Account) {
	tb.Helper()
	if err := m.HTTPServer.AccountService.CreateAccount(context.Background(), account); err != nil {
		tb.Fatal(err)
	}
}

func MustCloseMain(tb testing.TB, m *main.Main) {
	tb.Helper()
	if err := m.Close(); err != nil {
		tb.Fatal(err)
	}
}

func MustCreateContainer(tb testing.TB) *postgres.PostgresContainer {
	tb.Helper()
	container, err := postgres.Run(
		context.Background(),
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase("TestDB"),
		postgres.WithUsername("TestUser"),
		postgres.WithPassword("TestPassword"),
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
