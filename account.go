package hummingbird

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

type AccountFilter struct {
	ID *int
}

type AccountUpdate struct {
	Delta int
}

type AccountService interface {
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, id int, update *AccountUpdate) error
	FindAccountByID(ctx context.Context, id int) (*Account, error)
}
