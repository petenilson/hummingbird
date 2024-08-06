package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/petenilson/go-ledger"
)

func (s *Server) handleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ToAccountID   int `json:"to_account_id"`
		FromAccountID int `json:"from_account_id"`
		Amount        int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println(err)
		Error(w, r, &ledger.Error{Code: ledger.EINVALID, Message: "Invalid JSON"})
		return
	}

	transfer := ledger.NewTransfer(body.FromAccountID, body.ToAccountID, body.Amount, "")

	err := s.TransferService.CreateTransfer(r.Context(), transfer)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(transfer); err != nil {
		LogError(r, err)
		return
	}
}

type TransferService struct {
	Client *HTTPClient
}

func (c *LedgerClient) CreateTransfer(
	ctx context.Context,
	transfer *ledger.InterAccountTransfer,
) error {
	body, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	req, err := c.HTTPClient.newRequest("POST", "/transfers", bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		return parseResponseError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&transfer); err != nil {
		return err
	}

	return nil
}
