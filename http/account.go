package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/petenilson/go-ledger"
)

func (s *Server) handleGetAccountById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, &ledger.Error{Code: ledger.EINVALID, Message: "Invalid Account ID"})
		return
	}

	account, err := s.AccountService.FindAccountByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var account ledger.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		Error(w, r, &ledger.Error{Code: ledger.EINVALID, Message: "Invalid JSON"})
		return
	}

	err := s.AccountService.CreateAccount(r.Context(), &account)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		LogError(r, err)
		return
	}
}

type AccountService struct {
	Client *Client
}

func NewAccountService(client *Client) *AccountService {
	return &AccountService{Client: client}
}
func (as *AccountService) CreateAccount(ctx context.Context, account *ledger.Account) error {
	body, err := json.Marshal(account)
	if err != nil {
		return err
	}

	req, err := as.Client.newRequest("POST", "/accounts", bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		return parseResponseError(resp)
	}
	return nil
}

func (as *AccountService) FindAccountByID(ctx context.Context, id int) (*ledger.Account, error) {
	req, err := as.Client.newRequest("GET", fmt.Sprintf("/accounts/%d", id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, parseResponseError(resp)
	}

	var account ledger.Account
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}
	return &account, nil
}
