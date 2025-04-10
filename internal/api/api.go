package api

import (
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/logger"
	"github.com/aube/gophermart/internal/store"
)

type Server struct {
	logger *slog.Logger
	router *http.ServeMux
	store  store.Store
}

func NewRouter(store store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	s := &Server{
		logger: logger.New(),
		store:  store,
		router: mux,
	}

	s.configureRouter()

	return s.router
}

func (s *Server) configureRouter() {
	s.router.HandleFunc(`GET /user/balance`, http.HandlerFunc(s.UserBalance))
	s.router.HandleFunc(`GET /user/orders`, http.HandlerFunc(s.UserOrders))
	s.router.HandleFunc(`GET /user/withdrawals`, http.HandlerFunc(s.UserWithdrawals))
	s.router.HandleFunc(`POST /user/balance/withdraw`, http.HandlerFunc(s.UserBalanceWithdraw))
}
