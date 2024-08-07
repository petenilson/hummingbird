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
	m := MustRunMain(t)
	defer MustCloseMain(t, m)

	ctx := context.Background()

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{URL: m.HTTPServer.URL()},
	}

	// Create the accounts first to perform the transfer between.
	account_from, account_to := &hummingbird.Account{Name: "Account From"}, &hummingbird.Account{Name: "Account To"}
	MustCreateAccount(t, m, account_from)
	MustCreateAccount(t, m, account_to)

	// Create the transfer.
	transfer := hummingbird.NewTransfer(account_from.ID, account_to.ID, 10_000, "Testing")
	err := test_client.CreateTransfer(ctx, transfer)
	if err != nil {
		t.Fatal(err)
	} else if transfer.ID == 0 {
		t.Fatal()
	}

	// Find the entrys related to the transfer that we just created.
	entrys, count, err := test_client.FindEntrys(
		ctx, hummingbird.EntryFilter{AccountID: &account_from.ID},
	)
	if err != nil {
		t.Fatal(err)
	} else if count == 0 {
		t.Fatalf("Want count of 2, got %d", count)
	}

	// Check that entrys related to the transfer are what we expect.
	entry_want_debit := &hummingbird.Entry{
		AccountID: account_from.ID,
		CreatedAt: m.DB.Now(),
		Amount:    -10_000,
		Type:      "DEBIT",
	}
	entry_want_credit := &hummingbird.Entry{
		AccountID: account_to.ID,
		CreatedAt: m.DB.Now(),
		Amount:    10_000,
		Type:      "CREDIT",
	}
	opts := cmpopts.IgnoreFields(hummingbird.Entry{}, "ID")
	if eq := cmp.Equal(entry_want_debit, entrys[0], opts); eq != true {
		t.Fatalf(cmp.Diff(entry_want_debit, entrys[0]), opts)
	} else if eq := cmp.Equal(entry_want_credit, entrys[1], opts); eq != true {
		t.Fatalf(cmp.Diff(entry_want_credit, entrys[1]), opts)
	}
}
