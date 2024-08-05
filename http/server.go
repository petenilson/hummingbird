package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/petenilson/go-ledger"
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router

	Address string

	AccountService  ledger.AccountService
	TransferService ledger.TransferService
}

func NewServer(address string) *Server {
	s := &Server{
		server:  &http.Server{Addr: address},
		router:  mux.NewRouter(),
		Address: address,
	}
	// Set Not Found handler
	s.router.NotFoundHandler = http.HandlerFunc(handleNotFound)

	// Use middleware to set the default Content-type for all responses.
	s.router.Use(defaultContentTypeMiddleware)

	// Register routes here.
	// Accounts
	s.router.HandleFunc("/accounts/{id}", s.handleGetAccountById).Methods("GET")
	s.router.HandleFunc("/accounts", s.handleCreateAccount).Methods("POST")

	// Transfers
	s.router.HandleFunc("/transfers", s.handleCreateTransfer).Methods("POST")

	// Use the mux router as the handler.
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
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&ErrorResponse{Error: "Resourse Not Found."})
}

func defaultContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
