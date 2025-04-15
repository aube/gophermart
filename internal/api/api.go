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
	s.router.HandleFunc(`POST /user/register`, http.HandlerFunc(s.UserRegister))
	s.router.HandleFunc(`POST /user/login`, http.HandlerFunc(s.UserLogin))

	// Private
	s.router.HandleFunc(`GET /user/balance`, s.AuthMiddleware(s.UserBalance))
	s.router.HandleFunc(`GET /user/orders`, s.AuthMiddleware(s.UserOrders))
	s.router.HandleFunc(`GET /user/withdrawals`, s.AuthMiddleware(s.UserWithdrawals))
	s.router.HandleFunc(`POST /user/balance/withdraw`, s.AuthMiddleware(s.UserBalanceWithdraw))
	s.router.HandleFunc(`POST /user/orders`, s.AuthMiddleware(s.UploadUserOrders))

}
