package mocks

import (
	"context"
	"github.com/petenilson/hummingbird"
)

var _ hummingbird.TransferService = (*TransferService)(nil)

type TransferService struct {
	CreateTransferFn   func(context.Context, *hummingbird.Transfer) error
	FindTransferByIDFn func(context.Context, int) (*hummingbird.Transfer, error)
}

// FindTransferByID implements hummingbird.TransferService.
func (t *TransferService) FindTransferByID(
	ctx context.Context, id int,
) (*hummingbird.Transfer, error) {
	return t.FindTransferByIDFn(ctx, id)
}

// CreateTransfer implements hummingbird.TransferService.
func (t *TransferService) CreateTransfer(
	ctx context.Context, transfer *hummingbird.Transfer,
) error {
	return t.CreateTransferFn(ctx, transfer)
}
