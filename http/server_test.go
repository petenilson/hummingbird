package http_test

import (
	"io"
	"net/http"
	"testing"

	hbhttp "github.com/petenilson/hummingbird/http"
	"github.com/petenilson/hummingbird/mocks"
)

var TestAddress string = "http://localhost:8000"

type TestServer struct {
	*hbhttp.Server

	AccountService     mocks.AccountService
	TransactionService mocks.TransactionService
}

func MustOpenServer(tb testing.TB) *TestServer {
	tb.Helper()

	s := &TestServer{Server: hbhttp.NewServer("localhost:8000")}

	s.Server.AccountService = &s.AccountService
	s.Server.TransactionService = &s.TransactionService

	if err := s.Open(); err != nil {
		tb.Fatal(err)
	}
	return s
}

func MustCloseServer(tb testing.TB, s *TestServer) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}

func (s *TestServer) MustCreateNewRequest(
	tb testing.TB, method, url string, body io.Reader,
) *http.Request {
	tb.Helper()

	r, err := http.NewRequest(method, s.Address+url, body)
	if err != nil {
		tb.Fatal(err)
	}
	return r
}
