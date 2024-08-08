package http_test

import (
	"context"
	"testing"

	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/http"
)

func TestTransaction(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{TestAddress},
	}

	t.Run("CreateTransaction", func(t *testing.T) {
		s.TransactionService.CreateTransactionFn = func(context.Context, *hummingbird.Transaction) error {
			return nil
		}

		transaction := &hummingbird.Transaction{
			Description: "Test Transaction",
			Entrys: []*hummingbird.Entry{
				{
					AccountID: 1234,
					Amount:    -100,
					Type:      hummingbird.DEBIT,
				},
				{
					AccountID: 5678,
					Amount:    100,
					Type:      hummingbird.CREDIT,
				}},
		}

		err := test_client.CreateTransaction(context.Background(), transaction)
		if err != nil {
			t.Fatal(err)
		}
	})
}
