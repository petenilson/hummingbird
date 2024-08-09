package main_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/cmd/hummingbird"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var TestServer struct {
	URL string
}

var Services struct {
	AccountService hummingbird.AccountService
}

func TestMain(m *testing.M) {

	container, err := CreateContainer()
	if err != nil {
		log.Fatal(err)
	}
	dsn, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	application := main.NewMain()
	application.Config.DB.DSN = dsn
	application.Config.HTTP.Address = "localhost:8000"

	if err := application.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Set the URL of the web server so that test clients can connect.
	TestServer.URL = application.HTTPServer.URL()

	// Some tests require direct access to underlying services.
	Services.AccountService = application.HTTPServer.AccountService

	application.DB.Now = func() time.Time { return time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC) }

	// Run the all the tests
	exit_code := m.Run()

	// Clean up and shut down.
	if err := application.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exit_code)
}

func MustCreateAccount(
	tb testing.TB,
	as hummingbird.AccountService,
	account *hummingbird.Account,
) {
	tb.Helper()
	if err := as.CreateAccount(context.Background(), account); err != nil {
		tb.Fatal(err)
	}
}

func CreateContainer() (*postgres.PostgresContainer, error) {
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
		return nil, err
	}
	return container, nil
}
