package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/petenilson/hummingbird"
)

type TransactionService struct {
	db *DB
}

func NewTransactionService(db *DB) *TransactionService {
	return &TransactionService{db: db}
}

func (ts *TransactionService) CreateTransaction(ctx context.Context, transaction *hummingbird.Transaction) error {
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := createTransaction(ctx, tx, transaction); err != nil {
		return err
	}

	// TODO: This is slow. Bulk insert the entrys and transaction-entrys here
	for _, v := range transaction.Entrys {
		if err := createEntry(ctx, tx, v); err != nil {
			return err
		}
		if err := createTransactionEntry(
			ctx,
			tx,
			&hummingbird.TransactionEntry{
				CreatedAt:     v.CreatedAt,
				EntryID:       v.ID,
				TransactionID: transaction.ID,
			},
		); err != nil {
			return fmt.Errorf("CreateTransaction: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (ts *TransactionService) FindTransactionById(ctx context.Context, id int) (*hummingbird.Transaction, error) {
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if transaction, err := findTransactionById(ctx, tx, id); err != nil {
		return nil, err
	} else if err = attachEntrys(ctx, tx, transaction); err != nil {
		return nil, fmt.Errorf("FindTransactionById: %v", err)
	} else {
		return transaction, nil
	}
}

func createTransaction(ctx context.Context, tx *Tx, transaction *hummingbird.Transaction) error {
	transaction.CreatedAt = tx.asof
	// Insert row into database.
	err := tx.QueryRow(ctx, `
		INSERT INTO transactions (
			description,
			created_at
		)
		VALUES ($1, $2)
		RETURNING id
	`,
		transaction.Description,
		(*NullTime)(&transaction.CreatedAt),
	).Scan(&transaction.ID)
	if err != nil {
		return err
	}
	return nil
}

func findTransactionById(ctx context.Context, tx *Tx, id int) (*hummingbird.Transaction, error) {
	transactions, count, err := findTransactions(ctx, tx, hummingbird.TransactionFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	// TODO: Better error handling here
	if count == 0 {
		return nil, errors.New("No Transactions Found")
	}
	return transactions[0], nil
}

func findTransactions(ctx context.Context, tx *Tx, filter hummingbird.TransactionFilter) ([]*hummingbird.Transaction, int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	// Execue query with limiting WHERE clause and LIMIT/OFFSET injected.
	rows, err := tx.Query(ctx, `
		SELECT 
      id,
      created_at,
      description,
		  COUNT(*) OVER()
		FROM transactions
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	transactions := make([]*hummingbird.Transaction, 0)
	transactions_count := 0
	for rows.Next() {
		var transaction hummingbird.Transaction
		if err := rows.Scan(
			&transaction.ID,
			(*NullTime)(&transaction.CreatedAt),
			&transaction.Description,
			&transactions_count,
		); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return transactions, transactions_count, nil
}
