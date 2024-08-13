package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/petenilson/hummingbird"
)

func (s *Server) registerAccountRoutes(h huma.API) {
	// Get Account by id
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "get-account",
			Method:        http.MethodGet,
			Path:          "/accounts/{id}",
			Summary:       "Get Account",
			DefaultStatus: http.StatusOK,
		},
		s.handleGetAccountById,
	)

	// Create Account
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "create-account",
			Method:        http.MethodPost,
			Path:          "/accounts",
			Summary:       "Create Account",
			DefaultStatus: http.StatusCreated,
		},
		s.handleCreateAccount,
	)
}

func (s *Server) handleGetAccountById(
	ctx context.Context,
	req *struct {
		AccountID int `path:"id" maxLength:"30" example:"1" doc:"Account ID"`
	},
) (*Response[hummingbird.Account], error) {
	account, err := s.AccountService.FindAccountByID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	response := &Response[hummingbird.Account]{
		Body: account,
	}

	return response, nil
}

func (s *Server) handleCreateAccount(
	ctx context.Context,
	req *struct {
		Body struct {
			Name string `json:"name" example:"My Account Name" doc:"Account name"`
		}
	},
) (*Response[hummingbird.Account], error) {
	account := &hummingbird.Account{Name: req.Body.Name}
	err := s.AccountService.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	response := &Response[hummingbird.Account]{
		Body: account,
	}

	return response, nil
}

type AccountService struct {
	Client *HTTPClient
}

func NewAccountService(client *HTTPClient) *AccountService {
	return &AccountService{Client: client}
}

func (c *LedgerClient) CreateAccount(ctx context.Context, account *hummingbird.Account) error {
	r := struct {
		Name string `json:"name"`
	}{
		Name: account.Name,
	}
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := c.HTTPClient.newRequest("POST", "/accounts", bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		return parseResponseError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return err
	}

	return nil
}

func (c *LedgerClient) FindAccountByID(ctx context.Context, id int) (*hummingbird.Account, error) {
	req, err := c.HTTPClient.newRequest("GET", fmt.Sprintf("/accounts/%d", id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, parseResponseError(resp)
	}

	var account hummingbird.Account
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}
	return &account, nil
}
