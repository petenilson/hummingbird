package mocks

import (
	"context"

	"github.com/petenilson/hummingbird"
)

var _ hummingbird.AccountService = (*AccountService)(nil)

type AccountService struct {
	FindAccountByIDFn func(context.Context, int) (*hummingbird.Account, error)
	CreateAccountFn   func(context.Context, *hummingbird.Account) error
}

// CreateAccout implements hummingbird.AccountService.
func (s *AccountService) CreateAccount(ctx context.Context, account *hummingbird.Account) error {
	return s.CreateAccountFn(ctx, account)
}

// FindAccountByID implements hummingbird.AccountService.
func (s *AccountService) FindAccountByID(ctx context.Context, id int) (*hummingbird.Account, error) {
	return s.FindAccountByIDFn(ctx, id)
}

// UpdateAccount implements hummingbird.AccountService.
func (s *AccountService) UpdateAccount(
	ctx context.Context, account_id int, update *hummingbird.AccountUpdate,
) error {
	panic("unimplemented")
}
