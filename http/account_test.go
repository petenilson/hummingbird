package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/http"
)

func TestAccounts(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{"http://localhost:8000"},
	}

	account := &hummingbird.Account{
		ID:        999,
		Balance:   10_000,
		Name:      "Peter's Account",
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	s.AccountService.FindAccountByIDFn = func(context.Context, int) (*hummingbird.Account, error) {
		return account, nil
	}
	s.AccountService.CreateAccountFn = func(context.Context, *hummingbird.Account) error {
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
