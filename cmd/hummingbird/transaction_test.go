package main_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/http"
)

func TestTransactions(t *testing.T) {
	m := MustRunMain(t)
	defer MustCloseMain(t, m)

	ctx := context.Background()

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{URL: m.HTTPServer.URL()},
	}

	// Create the accounts first to perform the transfer between.
	account_from := &hummingbird.Account{Name: "Account From"}
	account_to := &hummingbird.Account{Name: "Account To"}
	MustCreateAccount(t, m, account_from)
	MustCreateAccount(t, m, account_to)

	// Create the transaction.
	transaction := &hummingbird.Transaction{
		Description: "Test Transaction",
		Entrys: []*hummingbird.Entry{
			{
				AccountID: account_from.ID,
				Amount:    -10_000,
				Type:      hummingbird.DEBIT,
			},
			{
				AccountID: account_to.ID,
				Amount:    10_000,
				Type:      hummingbird.CREDIT,
			}},
	}
	err := test_client.CreateTransaction(ctx, transaction)
	if err != nil {
		t.Fatal(err)
	} else if transaction.ID == 0 {
		t.Fatal()
	}

	// Find the entrys related to the transfer that we just created.
	if entrys, count, err := test_client.FindEntrys(
		ctx, hummingbird.EntryFilter{AccountID: &account_from.ID},
	); err != nil {
		t.Fatal(err)
	} else if count != 1 {
		t.Fatalf("Want count of 1, got %d", count)
	} else if diff := cmp.Diff(
		&hummingbird.Entry{
			AccountID: account_from.ID,
			CreatedAt: m.DB.Now(),
			Amount:    -10_000,
			Type:      hummingbird.DEBIT,
		},
		entrys[0],
		cmpopts.IgnoreFields(hummingbird.Entry{}, "ID"),
	); diff != "" {
		t.Fatalf("Want matching entrys but got: %s", diff)
	}

	if entrys, count, err := test_client.FindEntrys(
		ctx, hummingbird.EntryFilter{AccountID: &account_to.ID},
	); err != nil {
		t.Fatal(err)
	} else if count != 1 {
		t.Fatalf("Want count of 1, got %d", count)
	} else if diff := cmp.Diff(
		&hummingbird.Entry{
			AccountID: account_to.ID,
			CreatedAt: m.DB.Now(),
			Amount:    10_000,
			Type:      hummingbird.CREDIT,
		},
		entrys[0],
		cmpopts.IgnoreFields(hummingbird.Entry{}, "ID"),
	); diff != "" {
		t.Fatalf("Want matching entrys but got: %s", diff)
	}

}
