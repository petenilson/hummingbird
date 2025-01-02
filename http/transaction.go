package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/petenilson/hummingbird"
)

func (s *Server) registerTransactionRoutes(h huma.API) {
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "create-transaction",
			Method:        http.MethodPost,
			Path:          "/transactions",
			Summary:       "Create Transaction",
			DefaultStatus: http.StatusCreated,
		},
		s.handleCreateTransaction,
	)
	// Register Get Transaction
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "get-transaction",
			Method:        http.MethodGet,
			Path:          "/transactions/{transaction}",
			Summary:       "Get Transaction",
			DefaultStatus: http.StatusOK,
		},
		s.handleGetTransaction,
	)
	// Register List Transactions for Transfer
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "get-transactions-for-transfer",
			Method:        http.MethodGet,
			Path:          "/transfers/{transfer}/transactions",
			Summary:       "Get Transactions for Transfer",
			DefaultStatus: http.StatusOK,
		},
		s.handleGetTransaction,
	)
}

func (s *Server) handleGetTransaction(
	ctx context.Context,
	request *struct {
		Transaction int `path:"transaction"`
	},
) (*Response[hummingbird.Transaction], error) {
	transaction, err := s.TransactionService.FindTransactionByID(ctx, request.Transaction)
	if err != nil {
		return nil, err
	}

	response := &Response[hummingbird.Transaction]{
		Body: transaction,
	}

	return response, nil
}

func (s *Server) handleCreateTransaction(
	ctx context.Context,
	req *struct {
		Body struct {
			Description string
			Entrys      []*struct {
				AccountID int    `json:"account_id"`
				Amount    int    `json:"amount"`
				Type      string `json:"string" enum:"DEBIT,CREDIT"`
			}
		}
	},
) (*Response[hummingbird.Transaction], error) {
	transaction := &hummingbird.Transaction{
		Description: req.Body.Description,
	}
	for _, i := range req.Body.Entrys {
		e := &hummingbird.Entry{
			AccountID: i.AccountID,
			Amount:    i.Amount,
			Type:      hummingbird.EntryType(i.Type),
		}
		transaction.Entrys = append(transaction.Entrys, e)
	}

	err := s.TransactionService.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	response := &Response[hummingbird.Transaction]{
		Body: transaction,
	}

	return response, nil
}

func (c *LedgerClient) CreateTransaction(
	ctx context.Context,
	transaction *hummingbird.Transaction,
) error {
	body := struct {
		Description string
		Entrys      []struct {
			AccountID int    `json:"account_id"`
			Amount    int    `json:"amount"`
			Type      string `json:"string"`
		}
	}{
		Description: transaction.Description,
	}
	for _, e := range transaction.Entrys {
		body.Entrys = append(
			body.Entrys,
			struct {
				AccountID int    `json:"account_id"`
				Amount    int    `json:"amount"`
				Type      string `json:"string"`
			}{
				e.AccountID,
				e.Amount,
				string(e.Type),
			},
		)
	}
	body_bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := c.HTTPClient.newRequest("POST", "/transactions", bytes.NewReader(body_bytes))
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
