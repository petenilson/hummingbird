package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/petenilson/hummingbird/http"
	"github.com/petenilson/hummingbird/postgres"
)

const (
	DefaultConfigPath = "~./hummingbird.conf"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// Listen for termination of process.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	m := NewApplication()

	// Load config from config file.
	if err := m.ParseFlags(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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

type Application struct {
	Config Config

	DB *postgres.DB

	HTTPServer *http.Server
}

func NewApplication() *Application {
	return &Application{
		Config: Config{},
	}
}

func (app *Application) Run(ctx context.Context) error {
	app.HTTPServer = http.NewServer(app.Config.HTTP.Address)
	app.DB = postgres.NewDB(app.Config.DB.DSN)
	if err := app.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	app.HTTPServer.AccountService = postgres.NewAccountService(app.DB)
	app.HTTPServer.EntryService = postgres.NewEntryService(app.DB)
	app.HTTPServer.TransactionService = postgres.NewTransactionService(app.DB)

	if err := app.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

func (app *Application) Close() error {
	if app.HTTPServer != nil {
		if err := app.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) ParseFlags(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("hummingbird", flag.ContinueOnError)
	var config_path string
	fs.StringVar(&config_path, "config", DefaultConfigPath, "config path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	configPath, err := expand(config_path)
	if err != nil {
		return err
	}

	config, err := LoadConfig(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", config_path)
	} else if err != nil {
		return err
	}
	app.Config = config

	return nil
}

func expand(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return path, nil
	}

	u, err := user.Current()
	if err != nil {
		return path, err
	} else if u.HomeDir == "" {
		return path, fmt.Errorf("home directory unset")
	}

	if path == "~" {
		return u.HomeDir, nil
	}
	return filepath.Join(u.HomeDir, strings.TrimPrefix(path, "~"+string(os.PathSeparator))), nil
}

type Config struct {
	DB struct {
		DSN string `toml:"dsn"`
	} `toml:"db"`

	HTTP struct {
		Address string `toml:"address"`
		Domain  string `toml:"domain"`
	} `toml:"http"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	if buf, err := os.ReadFile(filename); err != nil {
		return config, err
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		return config, err
	}
	return config, nil
}
