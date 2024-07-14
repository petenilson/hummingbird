package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/petenilson/go-ledger"
)

type TransferService struct {
	db *DB
}

func NewTransferService(db *DB) *TransferService {
	return &TransferService{
		db: db,
	}
}

func (ts *TransferService) CreateTransfer(
	ctx context.Context, transfer *ledger.InterAccountTransfer,
) error {
	// Start Transaction
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock accounts rows
	if err := lockAccountRows(
		ctx, tx, []int{transfer.ToAccountID, transfer.FromAccountID},
	); err != nil {
		return err
	}

	// Update accounts balance
	if err := updateAccount(
		ctx,
		tx,
		transfer.FromAccountID,
		&ledger.AccountUpdate{
			Delta: -transfer.Amount,
		},
	); err != nil {
		return err
	} else if err := updateAccount(
		ctx,
		tx,
		transfer.ToAccountID,
		&ledger.AccountUpdate{
			Delta: transfer.Amount,
		},
	); err != nil {
		return err
	}

	// Create a new transaction
	if err := createTransaction(ctx, tx, transfer.Transaction); err != nil {
		return err
	}
	transfer.TransactionID = transfer.Transaction.ID

	// Create entries for transaction
	for _, v := range transfer.Transaction.Entrys {
		if err := createEntry(ctx, tx, v); err != nil {
			return err
		}
	}

	// Create the transfer
	if err := createTransfer(ctx, tx, transfer); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ts *TransferService) FindTransferById(
	ctx context.Context, id int,
) (*ledger.InterAccountTransfer, error) {
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	transfer, err := findTransferById(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	// TODO: this follows the foreign keys relationship
	// and introduces additional queries. Maybe give the option
	// to not follow relationships or do so in a more performant way.
	// if err := attachTransaction(ctx, tx, transfer); err != nil {
	// 	return nil, err
	// } else if err = attachEntrys(ctx, tx, transfer.Transaction); err != nil {
	// 	return nil, err
	// }

	return transfer, nil
}

func findTransferById(ctx context.Context, tx *Tx, id int) (*ledger.InterAccountTransfer, error) {
	transfers, count, err := findTransfers(ctx, tx, ledger.TransferFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	// TODO: Better error handling here
	if count == 0 {
		return nil, errors.New("No Transfer Found")
	}
	return transfers[0], nil
}

func findTransfers(
	ctx context.Context,
	tx *Tx,
	filter ledger.TransferFilter,
) (
	[]*ledger.InterAccountTransfer,
	int,
	error,
) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	rows, err := tx.Query(ctx, `
		SELECT 
      id,
      reason,
      from_account_id,
      to_account_id,
      created_at,
      amount,
    	transaction_id,
		  COUNT(*) OVER()
		FROM transfers
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	transfers := make([]*ledger.InterAccountTransfer, 0)
	transfer_count := 0
	for rows.Next() {
		var transfer ledger.InterAccountTransfer
		if err := rows.Scan(
			&transfer.ID,
			&transfer.Description,
			&transfer.FromAccountID,
			&transfer.ToAccountID,
			(*NullTime)(&transfer.CreatedAt),
			&transfer.Amount,
			&transfer.TransactionID,
			&transfer_count,
		); err != nil {
			return nil, 0, err
		}
		transfers = append(transfers, &transfer)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return transfers, transfer_count, nil
}

func attachTransaction(
	ctx context.Context, tx *Tx, transfer *ledger.InterAccountTransfer,
) error {
	if transaction, err := findTransactionById(
		ctx, tx, transfer.TransactionID,
	); err != nil {
		return err
	} else {
		transfer.Transaction = transaction
	}
	return nil
}

func createTransfer(ctx context.Context, tx *Tx, transfer *ledger.InterAccountTransfer) error {
	transfer.CreatedAt = tx.asof
	// Insert row into database.
	err := tx.QueryRow(ctx, `
		INSERT INTO transfers (
			from_account_id,
			to_account_id,
			amount,
			reason,
			created_at,
			transaction_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		transfer.FromAccountID,
		transfer.ToAccountID,
		transfer.Amount,
		transfer.Description,
		(*NullTime)(&transfer.CreatedAt),
		transfer.TransactionID,
	).Scan(&transfer.ID)
	if err != nil {
		return err
	}
	return nil
}
