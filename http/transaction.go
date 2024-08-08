package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/petenilson/hummingbird"
)

func (s *Server) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction hummingbird.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		Error(w, r, &hummingbird.Error{Code: hummingbird.EINVALID, Message: "Invalid JSON"})
		return
	}

	err := s.TransactionService.CreateTransaction(r.Context(), &transaction)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		LogError(r, err)
		return
	}
}

func (c *LedgerClient) CreateTransaction(
	ctx context.Context,
	transaction *hummingbird.Transaction,
) error {
	body, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	req, err := c.HTTPClient.newRequest("POST", "/transactions", bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		return parseResponseError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		return err
	}

	return nil
}
