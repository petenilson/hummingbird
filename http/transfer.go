package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/petenilson/go-ledger"
)

func (s *Server) handleCreateTransfer(w http.ResponseWriter, r *http.Request) {
	var transfer ledger.InterAccountTransfer
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		Error(w, r, &ledger.Error{Code: ledger.EINVALID, Message: "Invalid JSON"})
		return
	}

	err := s.TransferService.CreateTransfer(r.Context(), &transfer)
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
	Client *Client
}

func (ts *TransferService) CreateTransfer(
	ctx context.Context,
	transfer *ledger.InterAccountTransfer,
) error {
	body, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	req, err := ts.Client.newRequest("POST", "/transfers", bytes.NewReader(body))
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
