package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/petenilson/hummingbird"
	"github.com/petenilson/hummingbird/http"
)

func TestTransfer(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{TestAddress},
	}

	t.Run("FindTransfer", func(*testing.T) {
		s.TransferService.FindTransferByIDFn = func(
			context.Context,
			int,
		) (*hummingbird.Transfer, error) {
			return &hummingbird.Transfer{
				ID:            123,
				Description:   "Test Transfer",
				FromAccountID: 123,
				ToAccountID:   456,
				Amount:        100,
				CreatedAt:     time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
				Transaction:   &hummingbird.Transaction{},
				TransactionID: 0,
			}, nil
		}

		_, err := test_client.FindTransferByID(context.Background(), 123)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CreateTransfer", func(*testing.T) {
		s.TransferService.CreateTransferFn = func(context.Context, *hummingbird.Transfer) error { return nil }
		transfer := hummingbird.NewTransfer(1234, 5678, 10_000, "Testing")

		err := test_client.CreateTransfer(context.Background(), transfer)
		if err != nil {
			t.Fatal(err)
		}
	})

}
