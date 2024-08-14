package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/petenilson/hummingbird"
)

func (s *Server) registerTransferRoutes(h huma.API) {
	// Create a new Transfer
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "create-transfer",
			Method:        http.MethodPost,
			Path:          "/transfers",
			Summary:       "Create Transfer",
			DefaultStatus: http.StatusCreated,
		},
		s.handleCreateTransfer,
	)
}

func (s *Server) handleCreateTransfer(
	ctx context.Context,
	req *struct {
		Body struct {
			FromAccountID int `json:"from_account_id"`
			ToAccountID   int `json:"to_account_id"`
			Amount        int `json:"amount"`
		}
	},
) (*Response[hummingbird.Transfer], error) {
	transfer := &hummingbird.Transfer{
		FromAccountID: req.Body.FromAccountID,
		ToAccountID:   req.Body.ToAccountID,
		Amount:        req.Body.Amount,
	}

	if err := s.TransferService.CreateTransfer(ctx, transfer); err != nil {
		return nil, err
	}

	return &Response[hummingbird.Transfer]{
		Body: transfer,
	}, nil
}

func (c *LedgerClient) CreateTransfer(
	ctx context.Context,
	transfer *hummingbird.Transfer,
) error {
	req_body := struct {
		FromAccountID int `json:"from_account_id"`
		ToAccountID   int `json:"to_account_id"`
		Amount        int `json:"amount"`
	}{
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Amount:        transfer.Amount,
	}
	req_body_bytes, err := json.Marshal(req_body)
	if err != nil {
		return err
	}

	req, err := c.HTTPClient.newRequest("POST", "/transfers", bytes.NewReader(req_body_bytes))
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
