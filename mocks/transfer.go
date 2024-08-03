package mocks

import (
	"context"
	"github.com/petenilson/go-ledger"
)

var _ ledger.TransferService = (*TransferService)(nil)

type TransferService struct {
	CreateTransferFn func(context.Context, *ledger.InterAccountTransfer) error
}

// CreateTransfer implements ledger.TransferService.
func (t *TransferService) CreateTransfer(
	ctx context.Context, transfer *ledger.InterAccountTransfer,
) error {
	return t.CreateTransferFn(ctx, transfer)
}
