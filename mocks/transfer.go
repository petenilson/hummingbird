package mocks

import (
	"context"
	"github.com/petenilson/hummingbird"
)

var _ hummingbird.TransferService = (*TransferService)(nil)

type TransferService struct {
	CreateTransferFn func(context.Context, *hummingbird.Transfer) error
}

// CreateTransfer implements hummingbird.TransferService.
func (t *TransferService) CreateTransfer(
	ctx context.Context, transfer *hummingbird.Transfer,
) error {
	return t.CreateTransferFn(ctx, transfer)
}
