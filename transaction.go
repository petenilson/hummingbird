package ledger

import (
	"context"
	"time"
)

type EntryType string

const (
	DEBIT  EntryType = "DEBIT"
	CREDIT EntryType = "CREDIT"
)

// Transaction represents a complete financial event
// with at least two or more entrys. The entrys of a
// transaction should balance out.
type Transaction struct {
	ID          int
	Description string
	CreatedAt   time.Time
	Entrys      []*Entry
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	FindTransactionByID(ctx context.Context, id int) (*Transaction, error)
}

type TransactionFilter struct {
	ID *int
}

type Entry struct {
	ID            int
	AccountID     int
	TransactionID int
	CreatedAt     time.Time
	Amount        int
	Type          EntryType
}

type EntryFilter struct {
	TransactionID *int
}
