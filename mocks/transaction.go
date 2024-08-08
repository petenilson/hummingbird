package mocks

import (
	"context"
	"github.com/petenilson/hummingbird"
)

var _ hummingbird.TransactionService = (*TransactionService)(nil)

type TransactionService struct {
	CreateTransactionFn   func(context.Context, *hummingbird.Transaction) error
	FindTransactionByIDFn func(context.Context, int) (*hummingbird.Transaction, error)
}

// CreateTransaction implements hummingbird.TransactionService.
func (t *TransactionService) CreateTransaction(
	ctx context.Context,
	transaction *hummingbird.Transaction,
) error {
	return t.CreateTransactionFn(ctx, transaction)
}

// FindTransactionByID implements hummingbird.TransactionService.
func (t *TransactionService) FindTransactionByID(
	ctx context.Context,
	id int,
) (*hummingbird.Transaction, error) {
	return t.FindTransactionByIDFn(ctx, id)
}
