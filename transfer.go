package ledger

import "time"

// InterAccountTransfer represents a movement of funds with exactly one
// Transction associated with it.
type InterAccountTransfer struct {
	ID            int
	Description   string
	ToAccountID   int
	FromAccountID int
	Amount        int
	CreatedAt     time.Time
	Transaction   *Transaction
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
