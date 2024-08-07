package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/postgres"
)

func TestTransactionService_FindTransactionById(t *testing.T) {
	t.Run("OK", func(*testing.T) {
		cx := context.Background()
		// I want to create a new TX here
		// Functions being tested should use that transaction
		// to create nested transactions within.
		// I want to then rollbock the outermost transaction
		ts := postgres.NewTransactionService(DB)
		as := postgres.NewAccountService(DB)

		to_account := &hummingbird.Account{Name: "To Account"}
		from_account := &hummingbird.Account{Name: "From Account"}
		if err := as.CreateAccount(cx, to_account); err != nil {
			t.Fatal(err)
		} else if err := as.CreateAccount(cx, from_account); err != nil {
			t.Fatal(err)
		}

		tx := &hummingbird.Transaction{
			Description: "Test Transaction",
			Entrys: []*hummingbird.Entry{
				{
					AccountID: from_account.ID,
					Amount:    -100,
					Type:      hummingbird.DEBIT,
				},
				{
					AccountID: to_account.ID,
					Amount:    100,
					Type:      hummingbird.CREDIT,
				},
			},
		}
		if err := ts.CreateTransaction(cx, tx); err != nil {
			t.Fatal(err)
		}
		if result, err := ts.FindTransactionById(cx, tx.ID); err != nil {
			t.Fatal(err)
		} else if got, want := result.ID, tx.ID; got != want {
			t.Fatalf("ID=%v, want %v", got, want)
		} else if got, want := result.Description, tx.Description; got != want {
			t.Fatalf("Description =%v, want %v", got, want)
		} else if got, want := result.CreatedAt, tx.CreatedAt; got != want {
			t.Fatalf("CreatedAt=%v, want %v", got, want)
		} else if got, want := result.Entrys, tx.Entrys; cmp.Equal(got, want) != true {
			t.Fatalf(cmp.Diff(got, want))
		}
	})
}
