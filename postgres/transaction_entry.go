package postgres

import (
	"context"

	"github.com/petenilson/go-ledger"
)

type TransactionEntryService struct {
	db *DB
}

func NewTransactionEntryService(db *DB) *TransactionEntryService {
	return &TransactionEntryService{
		db: db,
	}
}

// FindEntryByTransactionID retreives entries in the context of a ledger.Transction
// and performs the lookup through the TransactionEntry table.
func (tas *TransactionEntryService) FindEntrsyByTransactionID(
	ctx context.Context, transaction_id int,
) ([]*ledger.Entry, int, error) {
	tx, err := tas.db.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)
	return findEntriesByTransactionID(ctx, tx, transaction_id)
}

func createTransactionEntry(
	ctx context.Context,
	tx *Tx,
	transaction_entry *ledger.TransactionEntry,
) error {
	transaction_entry.CreatedAt = tx.asof
	// Insert row into database.
	err := tx.QueryRow(ctx, `
		INSERT INTO transaction_entrys (
			transaction_id,
			entry_id,
			created_at
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`,
		transaction_entry.TransactionID,
		transaction_entry.EntryID,
		(*NullTime)(&transaction_entry.CreatedAt),
	).Scan(&transaction_entry.ID)
	if err != nil {
		return err
	}
	return nil
}

func findEntriesByTransactionID(
	ctx context.Context, tx *Tx, transaction_id int,
) ([]*ledger.Entry, int, error) {
	rows, err := tx.Query(ctx, `
		SELECT 
      e.id,
      e.account_id,
      e.created_at,
      e.amount,
    	e.type,
		  COUNT(*) OVER()
		FROM entrys e
		JOIN transaction_entrys te ON e.id = te.entry_id
		WHERE te.transaction_id = $1
		ORDER BY id ASC
		`,
		transaction_id,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entrys := make([]*ledger.Entry, 0)
	entry_count := 0

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
