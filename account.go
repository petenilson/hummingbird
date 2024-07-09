package ledger

import (
	"context"
	"time"
)

type Account struct {
	CreatedAt time.Time
	ID        int
	Balance   int
	Name      string
}

type AccountUpdate struct {
	Delta int
}

type AccountService interface {
	CreateAccout(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *AccountUpdate) (*Account, error)
	FindAccountByID(ctx context.Context, id int) (*Account, error)
}
