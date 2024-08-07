package postgres_test

import (
	"context"
	"testing"

	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/postgres"
)

func TestTransactionService_FindEntrysByTransactionID(t *testing.T) {
	t.Run("OK", func(*testing.T) {
		cx := context.Background()
		// tes := postgres.NewTransactionEntryService(TestDB)
		es := postgres.NewEntryService(DB)
		as := postgres.NewAccountService(DB)

		to_account := &hummingbird.Account{Name: "To Account"}
		fm_account := &hummingbird.Account{Name: "From Account"}
		if err := as.CreateAccount(cx, to_account); err != nil {
			t.Fatal(err)
		} else if err := as.CreateAccount(cx, fm_account); err != nil {
			t.Fatal(err)
		}

		entry_to_account := &hummingbird.Entry{
			Amount:    100,
			AccountID: to_account.ID,
			Type:      "CREDIT",
		}
		entry_fm_account := &hummingbird.Entry{
			Amount:    -100,
			AccountID: to_account.ID,
			Type:      "DEBIT",
		}
		if err := es.CreateEntry(cx, entry_to_account); err != nil {
			t.Fatal(err)
		} else if err := es.CreateEntry(cx, entry_fm_account); err != nil {
			t.Fatal(err)
		}

	})
}
