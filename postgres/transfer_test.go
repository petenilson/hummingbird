package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/petenilson/go-ledger"
	"github.com/petenilson/go-ledger/postgres"
)

func TestTransactionService_CreateTransfer(t *testing.T) {
	t.Run("OK", func(*testing.T) {
		cx := context.Background()
		db, clean_up := MustOpenDB(t)
		defer clean_up()
		ts := postgres.NewTransferService(db)
		as := postgres.NewAccountService(db)

		to_account := &ledger.Account{Name: "To Account"}
		from_account := &ledger.Account{Name: "From Account"}
		if err := as.CreateAccount(cx, to_account); err != nil {
			t.Fatal(err)
		} else if err := as.CreateAccount(cx, from_account); err != nil {
			t.Fatal(err)
		}

		transfer := ledger.NewTransfer(
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
		want := &ledger.InterAccountTransfer{
			ID:            transfer.ID,
			Description:   "Test Transfer",
			FromAccountID: from_account.ID,
			ToAccountID:   to_account.ID,
			Amount:        100,
			CreatedAt:     db.Now(),
			TransactionID: 1,
		}
		if eq := cmp.Equal(got, want); eq != true {
			t.Fatalf(cmp.Diff(got, want))
		}
	})
}
