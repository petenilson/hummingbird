package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/petenilson/go-ledger"
)

type AccountService struct {
	db *DB
}

func NewAccountService(db *DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (as *AccountService) CreateAccount(
	ctx context.Context,
	account *ledger.Account,
) error {
	tx, err := as.db.BeginTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	if err := createAccount(ctx, tx, account); err != nil {
		return err
	}
	return nil

}

func createAccount(ctx context.Context, tx *Tx, account *ledger.Account) error {
	account.CreatedAt = tx.asof
	// Insert row into database.
	err := tx.QueryRow(ctx, `
		INSERT INTO accounts (
			name,
			created_at,
			updated_at,
			balance
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		account.Name,
		(*NullTime)(&tx.asof),
		(*NullTime)(&tx.asof),
		0,
	).Scan(&account.ID)
	if err != nil {
		return err
	}
	return nil
}

func lockAccountRows(ctx context.Context, tx *Tx, account_ids []int) error {

	ids := make([]string, len(account_ids))
	args := make([]interface{}, len(account_ids))
	for i, id := range account_ids {
		ids[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	if _, err := tx.Exec(
		ctx,
		fmt.Sprintf(
			`SELECT *
			FROM accounts
			WHERE id IN (%s)
			FOR UPDATE`,
			strings.Join(ids, ", "),
		),
		args...); err != nil {
		return err
	}
	return nil
}

func updateAccount(
	ctx context.Context,
	tx *Tx,
	id int,
	update *ledger.AccountUpdate,
) error {
	// Execute update query.
	if _, err := tx.Exec(ctx, `
    UPDATE accounts 
    SET balance = balance + $1,
      updated_at = $2
    WHERE id = $3
	`,
		update.Delta,
		(*NullTime)(&tx.asof),
		id,
	); err != nil {
		return err
	}
	return nil
}
