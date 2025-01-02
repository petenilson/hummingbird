package hummingbird

import (
	"context"
	"time"
)

// Transfer represents a movement of funds between two accounts.
type Transfer struct {
	ID            int          `json:"id"`
	Description   string       `json:"description"`
	FromAccountID int          `json:"from_account_id"`
	ToAccountID   int          `json:"to_account_id"`
	Amount        int          `json:"amount"`
	CreatedAt     time.Time    `json:"created_at"`
	Transaction   *Transaction `json:"transaction"`
	TransactionID int          `json:"transaction_id"`
}

func NewTransfer(from_account_id, to_account_id, amount int, reason string) *Transfer {
	return &Transfer{
		Description:   reason,
		Amount:        amount,
		FromAccountID: from_account_id,
		ToAccountID:   to_account_id,
		Transaction: &Transaction{
			Entrys: []*Entry{
				{
					AccountID: from_account_id,
					Amount:    -amount,
					Type:      "DEBIT",
				},
				{
					AccountID: to_account_id,
					Amount:    amount,
					Type:      "CREDIT",
				},
			},
		},
	}
}

type TransferFilter struct {
	ID *int
}

type TransferService interface {
	CreateTransfer(ctx context.Context, transfer *Transfer) error
	FindTransferByID(ctx context.Context, id int) (*Transfer, error)
}
