package hummingbird

import (
	"context"
	"time"
)

// Transaction represents a complete financial event with at least two or more entrys.
// The entrys of a transaction should balance out.
type Transaction struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Entrys      []*Entry  `json:"entrys"`
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	FindTransactionByID(ctx context.Context, id int) (*Transaction, error)
}

type TransactionFilter struct {
	ID *int
}
