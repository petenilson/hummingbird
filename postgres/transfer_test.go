package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/postgres"
)

func TestTransactionService(t *testing.T) {
	t.Run("CreateTransfer", func(*testing.T) {
		ctx := context.Background()
		transfer_service := postgres.NewTransferService(DB)
		transaction_service := postgres.NewTransactionService(DB)
		account_service := postgres.NewAccountService(DB)

		to_account := &hummingbird.Account{Name: "To Account"}
		from_account := &hummingbird.Account{Name: "From Account"}
		if err := account_service.CreateAccount(ctx, to_account); err != nil {
			t.Fatal(err)
		} else if err := account_service.CreateAccount(ctx, from_account); err != nil {
			t.Fatal(err)
		}

		transfer := hummingbird.NewTransfer(from_account.ID, to_account.ID, 100, "Test Transfer")
		if err := transfer_service.CreateTransfer(ctx, transfer); err != nil {
			t.Fatal(err)
		}

		transfer_got, err := transfer_service.FindTransferByID(ctx, transfer.ID)
		if err != nil {
			t.Fatal(err)
		}

		// Test Transfer was created in DB
		transfer_want := &hummingbird.Transfer{
			ID:            transfer.ID,
			Description:   "Test Transfer",
			FromAccountID: from_account.ID,
			ToAccountID:   to_account.ID,
			Amount:        100,
			CreatedAt:     DB.Now(),
			TransactionID: transfer_got.TransactionID,
		}
		if eq := cmp.Equal(transfer_got, transfer_want); eq != true {
			t.Fatalf(cmp.Diff(transfer_got, transfer_want))
		}

		// Test Transaction was created for transfer
		transaction_got, err := transaction_service.FindTransactionByID(ctx, transfer.TransactionID)
		if err != nil {
			t.Fatal(err)
		}

		transaction_want := &hummingbird.Transaction{
			ID:          transfer.TransactionID,
			Description: "",
			CreatedAt:   DB.Now(),
			Entrys: []*hummingbird.Entry{
				{AccountID: transfer.FromAccountID, CreatedAt: DB.Now(), Amount: -100, Type: "DEBIT"},
				{AccountID: transfer.ToAccountID, CreatedAt: DB.Now(), Amount: 100, Type: "CREDIT"},
			},
		}
		if diff := cmp.Diff(
			transaction_want,
			transaction_got,
			cmpopts.IgnoreFields(hummingbird.Entry{}, "ID"),
		); diff != "" {
			t.Fatalf("Want %v got %v diff %s", transaction_want, transaction_got, diff)
		}
	})

}
