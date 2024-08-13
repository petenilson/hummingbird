package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
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

	h := humago.New(s.router, huma.DefaultConfig("Hummingbird", "1.0.0"))

	// Register Account Routes
	s.registerAccountRoutes(h)

	// Register Transaction Routes
	s.registerTransactionRoutes(h)

	// Register Entry Routes
	s.registerEntryRoutes(h)

	// Set Not Found handler
	s.router.HandleFunc("/", handleNotFound)

	// Use the http mux router as the handler.
	s.server.Handler = s.router

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

type Response[T any] struct {
	Body *T
}

func defaultContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
