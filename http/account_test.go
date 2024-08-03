package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/petenilson/go-ledger"
	"github.com/petenilson/go-ledger/http"
)

func TestAccounts(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	test_client := http.AccountService{
		Client: &http.Client{"http://localhost:8000"},
	}

	account := &ledger.Account{
		ID:        999,
		Balance:   10_000,
		Name:      "Peter's Account",
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	s.AccountService.FindAccountByIDFn = func(context.Context, int) (*ledger.Account, error) {
		return account, nil
	}
	s.AccountService.CreateAccountFn = func(context.Context, *ledger.Account) error {
		return nil
	}

	t.Run("Test can retrieve account by id.", func(t *testing.T) {

		_, err := test_client.FindAccountByID(context.Background(), 999)
		if err != nil {
			t.Error(err)
		}
	},
	)

	t.Run("Test can create account.", func(t *testing.T) {

		err := test_client.CreateAccount(context.Background(), account)
		if err != nil {
			t.Error(err)
		}
	},
	)

	t.Run("Test can create account.", func(t *testing.T) {

		err := test_client.CreateAccount(context.Background(), account)
		if err != nil {
			t.Error(err)
		}
	},
	)
}
