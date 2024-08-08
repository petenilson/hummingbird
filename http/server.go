package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/petenilson/hummingbird"
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *http.ServeMux

	Address string

	EntryService       hummingbird.EntryService
	AccountService     hummingbird.AccountService
	TransactionService hummingbird.TransactionService
}

func NewServer(address string) *Server {
	s := &Server{
		server:  &http.Server{Addr: address},
		router:  http.NewServeMux(),
		Address: address,
	}

	// Register Account Routes
	s.router.HandleFunc("GET /accounts/{id}", s.handleGetAccountById)
	s.router.HandleFunc("POST /accounts", s.handleCreateAccount)

	// Register Transaction Routes
	s.router.HandleFunc("POST /transactions", s.handleCreateTransaction)

	// Register Entry Routes
	s.router.HandleFunc("GET /entrys", s.handleListEntrys)

	// Set Not Found handler
	s.router.HandleFunc("/", handleNotFound)

	// Use the http mux router as the handler.
	s.server.Handler = defaultContentTypeMiddleware(s.router)

	return s
}

func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Address); err != nil {
		return err
	}

	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) Close() error {
	// Give the server time to finish serving active requests.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s", s.Address)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}
}

func defaultContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
