package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/petenilson/go-ledger"
)

func (s *Server) handleListEntrys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filter := &ledger.EntryFilter{}
	if value, ok := vars["account_id"]; ok {
		if id, err := strconv.Atoi(value); err != nil {
			Error(w, r, &ledger.Error{Code: ledger.EINVALID, Message: "Invalid account_id"})
		} else {
			filter.AccountID = &id
		}
	}

	entrys, _, err := s.EntryService.FindEntrys(r.Context(), *filter)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entrys); err != nil {
		LogError(r, err)
		return
	}
}

type EntryService struct {
	Client *HTTPClient
}

func NewEntryService(client *HTTPClient) *EntryService {
	return &EntryService{Client: client}
}

func (es *LedgerClient) FindEntrys(
	ctx context.Context, filter ledger.EntryFilter,
) ([]*ledger.Entry, int, error) {
	req, err := es.HTTPClient.newRequest(
		"GET", fmt.Sprintf("/entrys?account_id=%d", *filter.AccountID), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, 0, parseResponseError(resp)
	}

	var entrys []*ledger.Entry
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&entrys); err != nil {
		return nil, 0, err
	}
	return entrys, len(entrys), nil
}
