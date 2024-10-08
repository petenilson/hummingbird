package main_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/http"
)

func TestAccounts(t *testing.T) {

	ctx := context.Background()

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{URL: TestServer.URL},
	}

	t.Run("CreateAccount", func(t *testing.T) {
		// Create Account through our Test Client.
		account_want := &hummingbird.Account{Name: "Test Account"}
		if err := test_client.CreateAccount(ctx, account_want); err != nil {
			t.Fatal(err)
		}

		// Retrieve Account through the Test Client.
		account_got, err := test_client.FindAccountByID(ctx, account_want.ID)
		if err != nil {
			t.Fatal(err)
		}

		// Ensure both are the same.
		if eq := cmp.Equal(account_want, account_got); eq != true {
			t.Fatalf(cmp.Diff(account_want, account_got))
		}
	},
	)

}
