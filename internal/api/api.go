package api

import (
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/logger"
	"github.com/aube/gophermart/internal/middlewares"
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
	s.router.HandleFunc(`GET /user/balance`, middlewares.AuthMiddleware(s.UserBalance))
	s.router.HandleFunc(`GET /user/orders`, http.HandlerFunc(s.UserOrders))
	s.router.HandleFunc(`GET /user/withdrawals`, http.HandlerFunc(s.UserWithdrawals))
	s.router.HandleFunc(`POST /user/balance/withdraw`, http.HandlerFunc(s.UserBalanceWithdraw))
	s.router.HandleFunc(`POST /user/register`, http.HandlerFunc(s.UserRegister))
	s.router.HandleFunc(`POST /user/login`, http.HandlerFunc(s.UserLogin))
	s.router.HandleFunc(`POST /user/orders`, http.HandlerFunc(s.UploadUserOrders))
}
