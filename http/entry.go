package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
			Path:          "/entrys",
			Summary:       "List Entrys",
			DefaultStatus: http.StatusOK,
		},
		s.handleListEntrys,
	)
}

func (s *Server) handleListEntrys(
	ctx context.Context,
	req *struct {
		AccountID int `query:"account_id" maxLength:"30" example:"1" doc:"Account ID"`
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
		fmt.Println(resp)
		return nil, 0, parseResponseError(resp)
	}

	var entrys []*hummingbird.Entry
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&entrys); err != nil {
		return nil, 0, err
	}
	return entrys, len(entrys), nil
}
