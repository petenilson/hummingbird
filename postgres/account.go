package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/petenilson/hummingbird"
)

var _ hummingbird.AccountService = (*AccountService)(nil)

type AccountService struct {
	db *DB
}

func NewAccountService(db *DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

// UpdateAccount implements ledger.AccountService.
func (as *AccountService) UpdateAccount(
	ctx context.Context, account_id int, update *hummingbird.AccountUpdate,
) error {
	tx, err := as.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	if err := updateAccount(ctx, tx, account_id, update); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// CreateAccount implements ledger.AccountService.
func (as *AccountService) CreateAccount(
	ctx context.Context,
	account *hummingbird.Account,
) error {
	tx, err := as.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}
	if err := createAccount(ctx, tx, account); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FindAccountByID implements ledger.AccountService.
func (as *AccountService) FindAccountByID(
	ctx context.Context, account_id int,
) (*hummingbird.Account, error) {
	tx, err := as.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	return findAccountById(ctx, tx, account_id)
}

func findAccountById(
	ctx context.Context, tx *Tx, account_id int,
) (*hummingbird.Account, error) {
	accounts, count, err := findAccounts(
		ctx,
		tx,
		&hummingbird.AccountFilter{
			ID: &account_id,
		})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, &hummingbird.Error{Code: hummingbird.ENOTFOUND, Message: "Account Not Found"}
	} else {
		return accounts[0], nil
	}
}

func findAccounts(
	ctx context.Context, tx *Tx, filter *hummingbird.AccountFilter,
) ([]*hummingbird.Account, int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	// Execue query with limiting WHERE clause and LIMIT/OFFSET injected.
	rows, err := tx.Query(ctx, `
		SELECT 
      id,
    	balance,
    	name,
      created_at,
      updated_at,
		  COUNT(*) OVER()
		FROM accounts
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	accounts := make([]*hummingbird.Account, 0)
	count := 0
	for rows.Next() {
		var account hummingbird.Account
		if err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.Name,
			(*NullTime)(&account.CreatedAt),
			(*NullTime)(&account.UpdatedAt),
			&count,
		); err != nil {
			return nil, 0, err
		}
		accounts = append(accounts, &account)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return accounts, count, nil
}

func createAccount(ctx context.Context, tx *Tx, account *hummingbird.Account) error {
	account.CreatedAt = tx.asof
	account.UpdatedAt = tx.asof
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

func updateAccount(
	ctx context.Context,
	tx *Tx,
	id int,
	update *hummingbird.AccountUpdate,
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
		args...,
	); err != nil {
		return err
	}
	return nil
}
