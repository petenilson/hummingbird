package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/postgres"
)

func TestTransactionService_FindAccounts(t *testing.T) {
	t.Run("OK", func(*testing.T) {
		cx := context.Background()

		as := postgres.NewAccountService(DB)

		to_account := hummingbird.NewAccount("To Account")
		if err := as.CreateAccount(cx, to_account); err != nil {
			t.Fatal(err)
		}
		fm_account := hummingbird.NewAccount("From Account")
		if err := as.CreateAccount(cx, fm_account); err != nil {
			t.Fatal(err)
		}

		got_to_account, err := as.FindAccountByID(cx, to_account.ID)
		if err != nil {
			t.Fatal(err)
		}
		got_fm_account, err := as.FindAccountByID(cx, fm_account.ID)
		if err != nil {
			t.Fatal(err)
		}

		want_to_account := &hummingbird.Account{
			ID:        to_account.ID,
			Balance:   0,
			Name:      "To Account",
			CreatedAt: DB.Now(),
			UpdatedAt: DB.Now(),
		}
		if eq := cmp.Equal(got_to_account, want_to_account); eq != true {
			t.Fatalf(cmp.Diff(got_to_account, want_to_account))
		}

		want_fm_account := &hummingbird.Account{
			ID:        fm_account.ID,
			Balance:   0,
			Name:      "From Account",
			CreatedAt: DB.Now(),
			UpdatedAt: DB.Now(),
		}
		if eq := cmp.Equal(got_fm_account, want_fm_account); eq != true {
			t.Fatalf(cmp.Diff(got_fm_account, want_fm_account))
		}
	})
}
