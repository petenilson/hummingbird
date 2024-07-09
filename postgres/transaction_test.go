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

		tx := &ledger.Transaction{
			Description: "Test Transaction",
			Entrys: []*ledger.Entry{
				{
					Amount: -100,
					Type:   ledger.DEBIT,
				},
				{
					Amount: 100,
					Type:   ledger.CREDIT,
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
