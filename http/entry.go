package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/petenilson/hummingbird"
)

func (s *Server) registerEntryRoutes(h huma.API) {
	// Get Entrys for an Account
	huma.Register(
		h,
		huma.Operation{
			OperationID:   "list-entrys",
			Method:        http.MethodGet,
			Path:          "/accounts/{id}/entrys",
			Summary:       "List Entrys",
			DefaultStatus: http.StatusOK,
		},
		s.handleGetEntrysForAccountID,
	)
}

func (s *Server) handleGetEntrysForAccountID(
	ctx context.Context,
	req *struct {
		AccountID int `path:"id" maxLength:"30" example:"1" doc:"Account ID"`
	},
) (*Response[[]*hummingbird.Entry], error) {
	entrys, _, err := s.EntryService.FindEntrys(
		ctx,
		hummingbird.EntryFilter{
			AccountID: &req.AccountID,
		},
	)
	if err != nil {
		return nil, err
	}

	response := &Response[[]*hummingbird.Entry]{
		Body: &entrys,
	}

	return response, nil
}
func (s *Server) handleListEntrysOld(w http.ResponseWriter, r *http.Request) {
	filter := &hummingbird.EntryFilter{}
	if account_id := r.URL.Query().Get("account_id"); account_id != "" {
		if value, err := strconv.Atoi(account_id); err != nil {
			Error(w, r, &hummingbird.Error{Code: hummingbird.EINVALID, Message: "Invalid account_id"})
		} else {
			filter.AccountID = &value
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
	ctx context.Context, filter hummingbird.EntryFilter,
) ([]*hummingbird.Entry, int, error) {
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

	var entrys []*hummingbird.Entry
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&entrys); err != nil {
		return nil, 0, err
	}
	return entrys, len(entrys), nil
}
