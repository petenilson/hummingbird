package ledger

import (
	"context"
	"time"
)

// InterAccountTransfer represents a movement of funds between two accounts.
type InterAccountTransfer struct {
	ID            int          `json:"id"`
	Description   string       `json:"description"`
	ToAccountID   int          `json:"to_account_id"`
	FromAccountID int          `json:"from_account_id"`
	Amount        int          `json:"amount"`
	CreatedAt     time.Time    `json:"created_at"`
	Transaction   *Transaction `json:"transaction"`
	TransactionID int          `json:"transaction_id"`
}

func NewTransfer(
	from_account_id,
	to_account_id,
	amount int,
	reason string,
) *InterAccountTransfer {
	return &InterAccountTransfer{
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
	CreateTransfer(context.Context, *InterAccountTransfer) error
}
