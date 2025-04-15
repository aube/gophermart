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
	// Public
	s.router.HandleFunc(`POST /api/user/register`, http.HandlerFunc(s.UserRegister))
	s.router.HandleFunc(`POST /api/user/login`, http.HandlerFunc(s.UserLogin))

	// Private
	s.router.HandleFunc(`GET /api/user/orders`, s.AuthMiddleware(s.UserOrders))
	s.router.HandleFunc(`GET /api/user/balance`, s.AuthMiddleware(s.UserBalance))
	s.router.HandleFunc(`POST /api/user/orders`, s.AuthMiddleware(s.UploadUserOrders))
	s.router.HandleFunc(`GET /api/user/withdrawals`, s.AuthMiddleware(s.UserWithdrawals))
	s.router.HandleFunc(`POST /api/user/balance/withdraw`, s.AuthMiddleware(s.UserBalanceWithdraw))
}
