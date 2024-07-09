package postgres

import (
	"context"
	"database/sql/driver"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	conn_pool *pgxpool.Pool
	ctx       context.Context // background context
	dsn       string
	cancel    func() // cancel background context
	Now       func() time.Time
}

func NewDB(dsn string) *DB {
	db := &DB{
		dsn: dsn,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() error {
	if conn_pool, err := pgxpool.New(db.ctx, db.dsn); err != nil {
		return err
	} else {
		db.conn_pool = conn_pool
	}
	if err := db.migrate(); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.conn_pool != nil {
		db.conn_pool.Close()
	}
	return nil
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.conn_pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:   tx,
		db:   db,
		asof: db.Now().UTC().Truncate(time.Second),
	}, nil
}

func (db *DB) migrateFile(name string) error {
	tx, err := db.conn_pool.Begin(db.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(db.ctx)

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(db.ctx, `SELECT COUNT(*) FROM migrations WHERE name = $1`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationsFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(db.ctx, string(buf)); err != nil {
		return err
	}

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(db.ctx, `INSERT INTO migrations (name) VALUES ($1)`, name); err != nil {
		return err
	}

	return tx.Commit(db.ctx)
}
func (db *DB) migrate() error {
	if _, err := db.conn_pool.Exec(
		db.ctx,
		`CREATE TABLE IF NOT EXISTS migrations (
			name TEXT PRIMARY KEY
		);`,
	); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	return nil
}

// NullTime represents a helper wrapper for time.Time. It automatically converts
// time fields to/from RFC 3339 format. Also supports NULL for zero time.
type NullTime time.Time

// Scan reads a time value from the database.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	} else if value, ok := value.(string); ok {
		*(*time.Time)(n), _ = time.Parse(time.RFC3339, value)
		return nil
	}
	return fmt.Errorf("NullTime: cannot scan to time.Time: %T", value)
}

// Value formats a time value for the database.
func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).UTC().Format(time.RFC3339), nil
}

// Tx wraps the PGX Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	pgx.Tx
	db   *DB
	asof time.Time
}
