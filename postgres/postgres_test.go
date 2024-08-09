package postgres_test

import (
	"log"
	"os"
	"testing"
	"time"

	"context"

	"github.com/petenilson/hummingbird/postgres"

	"github.com/testcontainers/testcontainers-go"
	container_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var DB *postgres.DB

func TestMain(m *testing.M) {

	// Open new Test Container
	container, err := CreateContainer()
	if err != nil {
		log.Fatal(err)
	}

	// Get the connection string from container
	dsn, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Make the DB avaiable for all test cases and then open the DB.
	DB = postgres.NewDB(dsn)
	if err := DB.Open(); err != nil {
		log.Fatal(err)
	}

	// Mock the DB.Now so all entities created in the DB should be equal.
	DB.Now = func() time.Time {
		return time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
	}

	// Run the all the tests
	exit_code := m.Run()

	// Clean up by closing the test database and write and write the status code to STDOUT
	DB.Close()
	os.Exit(exit_code)
}

func CreateContainer() (*container_postgres.PostgresContainer, error) {
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
		return nil, err
	}
	return container, nil
}
