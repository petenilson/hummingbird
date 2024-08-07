package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/petenilson/hummingbird"
)

// Ledger Client Provides an API for the Ledger over HTTP.
type LedgerClient struct {
	HTTPClient *HTTPClient
}

type HTTPClient struct {
	URL string
}

func NewClient(u string) *HTTPClient {
	return &HTTPClient{URL: u}
}

func (c *HTTPClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.URL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	return req, nil
}

func Error(w http.ResponseWriter, r *http.Request, e error) {
	code, message := hummingbird.ErrorCode(e), hummingbird.ErrorMessage(e)

	if code == hummingbird.EINTERNAL {
		LogError(r, e)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func parseResponseError(resp *http.Response) error {
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errorResponse ErrorResponse
	if err := json.Unmarshal(buf, &errorResponse); err != nil {
		return err
	}

	return &hummingbird.Error{
		Message: errorResponse.Error,
		Code:    FromErrorStatusCode(resp.StatusCode),
	}
}

func LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

var codes = map[string]int{
	hummingbird.ECONFLICT:       http.StatusConflict,
	hummingbird.EINVALID:        http.StatusBadRequest,
	hummingbird.ENOTFOUND:       http.StatusNotFound,
	hummingbird.ENOTIMPLEMENTED: http.StatusNotImplemented,
	hummingbird.EUNAUTHORIZED:   http.StatusUnauthorized,
	hummingbird.EINTERNAL:       http.StatusInternalServerError,
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return hummingbird.EINTERNAL
}
