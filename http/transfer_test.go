package http_test

import (
	"context"
	"testing"

	"github.com/petenilson/go-ledger"
	"github.com/petenilson/go-ledger/http"
)

func TestCreateTransfer(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	test_client := http.LedgerClient{
		HTTPClient: &http.HTTPClient{TestAddress},
	}

	s.TransferService.CreateTransferFn = func(context.Context, *ledger.InterAccountTransfer) error { return nil }

	transfer := ledger.NewTransfer(1234, 5678, 10_000, "Testing")

	err := test_client.CreateTransfer(context.Background(), transfer)
	if err != nil {
		t.Fatal(err)
	}

}
