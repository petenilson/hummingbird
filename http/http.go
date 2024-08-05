package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/petenilson/go-ledger"
)

type Client struct {
	URL string
}

func NewClient(u string) *Client {
	return &Client{URL: u}
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.URL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	return req, nil
}

func Error(w http.ResponseWriter, r *http.Request, e error) {
	code, message := ledger.ErrorCode(e), ledger.ErrorMessage(e)

	if code == ledger.EINTERNAL {
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

	return &ledger.Error{
		Message: errorResponse.Error,
		Code:    FromErrorStatusCode(resp.StatusCode),
	}
}

func LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

var codes = map[string]int{
	ledger.ECONFLICT:       http.StatusConflict,
	ledger.EINVALID:        http.StatusBadRequest,
	ledger.ENOTFOUND:       http.StatusNotFound,
	ledger.ENOTIMPLEMENTED: http.StatusNotImplemented,
	ledger.EUNAUTHORIZED:   http.StatusUnauthorized,
	ledger.EINTERNAL:       http.StatusInternalServerError,
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
	return ledger.EINTERNAL
}
