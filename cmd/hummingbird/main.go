package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/petenilson/hummingbird/http"
	"github.com/petenilson/hummingbird/postgres"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	m := NewMain()
	m.Config.DB.DSN = "postgresql://ledger_user:ledger_password@localhost:5432/ledger_db"
	m.Config.HTTP.Address = "localhost:8000"

	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done()

	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	Config     Config
	ConfigPath string

	DB *postgres.DB

	HTTPServer *http.Server
}

func NewMain() *Main {
	config := DefaultConfig()
	return &Main{
		Config: config,
	}
}

func (m *Main) Run(ctx context.Context) error {
	m.HTTPServer = http.NewServer(m.Config.HTTP.Address)
	m.DB = postgres.NewDB(m.Config.DB.DSN)
	if err := m.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	m.HTTPServer.AccountService = postgres.NewAccountService(m.DB)
	m.HTTPServer.EntryService = postgres.NewEntryService(m.DB)
	m.HTTPServer.TransactionService = postgres.NewTransactionService(m.DB)

	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	DB struct {
		DSN string
	}

	HTTP struct {
		Address  string
		Domain   string
		HashKey  string
		BlockKey string
	}
}

func DefaultConfig() Config {
	var config Config
	config.DB.DSN = ""
	return config
}
