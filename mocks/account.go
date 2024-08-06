package mocks

import (
	"context"

	"github.com/petenilson/go-ledger"
)

var _ ledger.AccountService = (*AccountService)(nil)

type AccountService struct {
	FindAccountByIDFn func(context.Context, int) (*ledger.Account, error)
	CreateAccountFn   func(context.Context, *ledger.Account) error
}

// CreateAccout implements ledger.AccountService.
func (s *AccountService) CreateAccount(ctx context.Context, account *ledger.Account) error {
	return s.CreateAccountFn(ctx, account)
}

// FindAccountByID implements ledger.AccountService.
func (s *AccountService) FindAccountByID(ctx context.Context, id int) (*ledger.Account, error) {
	return s.FindAccountByIDFn(ctx, id)
}

// UpdateAccount implements ledger.AccountService.
func (s *AccountService) UpdateAccount(
	ctx context.Context, account_id int, update *ledger.AccountUpdate,
) error {
	panic("unimplemented")
}
