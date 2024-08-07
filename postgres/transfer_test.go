package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/postgres"
)

func TestTransactionService_CreateTransfer(t *testing.T) {
	t.Run("OK", func(*testing.T) {
		cx := context.Background()
		ts := postgres.NewTransferService(DB)
		as := postgres.NewAccountService(DB)

		to_account := &hummingbird.Account{Name: "To Account"}
		from_account := &hummingbird.Account{Name: "From Account"}
		if err := as.CreateAccount(cx, to_account); err != nil {
			t.Fatal(err)
		} else if err := as.CreateAccount(cx, from_account); err != nil {
			t.Fatal(err)
		}

		transfer := hummingbird.NewTransfer(
			from_account.ID,
			to_account.ID,
			100,
			"Test Transfer",
		)
		if err := ts.CreateTransfer(cx, transfer); err != nil {
			t.Fatal(err)
		}

		got, err := ts.FindTransferById(cx, transfer.ID)
		if err != nil {
			t.Fatal(err)
		}
		// Not testing for attached Transactions yet.
		want := &hummingbird.InterAccountTransfer{
			ID:            transfer.ID,
			Description:   "Test Transfer",
			FromAccountID: from_account.ID,
			ToAccountID:   to_account.ID,
			Amount:        100,
			CreatedAt:     DB.Now(),
			TransactionID: got.TransactionID,
		}
		if eq := cmp.Equal(got, want); eq != true {
			t.Fatalf(cmp.Diff(got, want))
		}
	})
}
