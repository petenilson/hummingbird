package http

import (
	"context"
	"io"
	"net/http"
)

// Client represents an HTTP client.
type Client struct {
	URL string
}

func NewClient(u string) *Client {
	return &Client{URL: u}
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.URL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	return req, nil
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

type ErrorResponse struct {
	Error string `json:"error"`
}
