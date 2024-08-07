package mocks

import (
	"context"
	"github.com/petenilson/hummingbird"
)

var _ hummingbird.TransferService = (*TransferService)(nil)

type TransferService struct {
	CreateTransferFn func(context.Context, *hummingbird.InterAccountTransfer) error
}

// CreateTransfer implements hummingbird.TransferService.
func (t *TransferService) CreateTransfer(
	ctx context.Context, transfer *hummingbird.InterAccountTransfer,
) error {
	return t.CreateTransferFn(ctx, transfer)
}
