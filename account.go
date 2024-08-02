package ledger

import (
	"context"
	"time"
)

type Account struct {
	ID        int
	Balance   int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(name string) *Account {
	return &Account{
		Balance: 0,
		Name:    name,
	}
}

type AccountFilter struct {
	ID *int
}

type AccountUpdate struct {
	Delta int
}

type AccountService interface {
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *AccountUpdate) (*Account, error)
	FindAccountByID(ctx context.Context, id int) (*Account, error)
}
