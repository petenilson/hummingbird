package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/petenilson/go-ledger"
	"github.com/petenilson/go-ledger/postgres"
)

func TestTransactionService_FindTransactionById(t *testing.T) {
	t.Run("OK", func(*testing.T) {
		cx := context.Background()
		db, clean_up := MustOpenDB(t)
		defer clean_up()
		ts := postgres.NewTransactionService(db)
		as := postgres.NewAccountService(db)

		to_account := &ledger.Account{Name: "To Account"}
		from_account := &ledger.Account{Name: "From Account"}
		if err := as.CreateAccount(cx, to_account); err != nil {
			t.Fatal(err)
		} else if err := as.CreateAccount(cx, from_account); err != nil {
			t.Fatal(err)
		}

		tx := &ledger.Transaction{
			Description: "Test Transaction",
			Entrys: []*ledger.Entry{
				{
					AccountID: from_account.ID,
					Amount:    -100,
					Type:      ledger.DEBIT,
				},
				{
					AccountID: to_account.ID,
					Amount:    100,
					Type:      ledger.CREDIT,
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
