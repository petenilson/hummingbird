package main_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/http"
)

func TestTransfers(t *testing.T) {

	ctx := context.Background()

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{URL: TestServer.URL},
	}

	// Create the accounts first to perform the transfer between.
	account_from, account_to := &hummingbird.Account{Name: "Account From"}, &hummingbird.Account{Name: "Account To"}
	MustCreateAccount(t, Services.AccountService, account_from)
	MustCreateAccount(t, Services.AccountService, account_to)

	// Create the transfer.
	transfer := hummingbird.NewTransfer(account_from.ID, account_to.ID, 10_000, "Testing")
	err := test_client.CreateTransfer(ctx, transfer)
	if err != nil {
		t.Fatal(err)
	} else if transfer.ID == 0 {
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
			Amount:    -10_000,
			Type:      hummingbird.DEBIT,
		},
		entrys[0],
		cmpopts.IgnoreFields(hummingbird.Entry{}, "ID", "CreatedAt"),
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
			Amount:    10_000,
			Type:      hummingbird.CREDIT,
		},
		entrys[0],
		cmpopts.IgnoreFields(hummingbird.Entry{}, "ID", "CreatedAt"),
	); diff != "" {
		t.Fatalf("Want matching entrys but got: %s", diff)
	}
}
