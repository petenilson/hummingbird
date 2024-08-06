package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/petenilson/go-ledger"
)

var _ ledger.EntryService = (*EntryService)(nil)

type EntryService struct {
	db *DB
}

func NewEntryService(db *DB) *EntryService {
	return &EntryService{
		db: db,
	}
}

// FindEntrys implements ledger.EntryService.
func (es *EntryService) FindEntrys(
	ctx context.Context, filter ledger.EntryFilter,
) ([]*ledger.Entry, int, error) {
	tx, err := es.db.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)
	return findEntrys(ctx, tx, &filter)
}

// CreateEntry implements ledger.EntryService.
func (es *EntryService) CreateEntry(ctx context.Context, entry *ledger.Entry) error {
	tx, err := es.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createEntry(ctx, tx, entry); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func createEntry(ctx context.Context, tx *Tx, entry *ledger.Entry) error {
	entry.CreatedAt = tx.asof
	// Insert row into database.
	err := tx.QueryRow(ctx, `
		INSERT INTO entrys (
			account_id,
			created_at,
			amount,
			type
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		entry.AccountID,
		(*NullTime)(&entry.CreatedAt),
		entry.Amount,
		entry.Type,
	).Scan(&entry.ID)
	if err != nil {
		return err
	}
	return nil
}

// attachEntrys get's the entries associated with the transaction
// from the DB and adds them to the ledger.Transaction instance.
func attachEntrys(
	ctx context.Context, tx *Tx, transaction *ledger.Transaction,
) error {
	if entries, _, err := findEntriesByTransactionID(
		ctx,
		tx,
		transaction.ID,
	); err != nil {
		return err
	} else {
		fmt.Printf("%+v", entries)
		transaction.Entrys = entries
	}
	return nil
}

func findEntrys(
	ctx context.Context, tx *Tx, filter *ledger.EntryFilter,
) ([]*ledger.Entry, int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.AccountID; v != nil {
		where, args = append(where, "account_id = $1"), append(args, *v)
	}
	rows, err := tx.Query(ctx, `
		SELECT 
      id,
    	account_id,
      created_at,
    	amount,
    	type,
		  COUNT(*) OVER()
		FROM entrys
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	entrys := make([]*ledger.Entry, 0)
	var entry_count int
	for rows.Next() {
		var entry ledger.Entry
		if err := rows.Scan(
			&entry.ID,
			&entry.AccountID,
			(*NullTime)(&entry.CreatedAt),
			&entry.Amount,
			&entry.Type,
			&entry_count,
		); err != nil {
			return nil, 0, err
		}
		entrys = append(entrys, &entry)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return entrys, entry_count, nil
}
