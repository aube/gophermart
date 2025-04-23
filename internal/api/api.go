package api

import (
	"net/http"

	"github.com/aube/gophermart/internal/logger"
	"github.com/aube/gophermart/internal/store"
)

func NewRouter(store store.Store) *http.ServeMux {
	mux := http.NewServeMux()
	logger := logger.New()
	AuthMiddleware := NewAuthMiddleware(store.ActiveUser, logger)

	// Public
	mux.HandleFunc(`POST /api/user/register`, NewUserRegisterHandler(store.User, store.ActiveUser, logger))
	mux.HandleFunc(`POST /api/user/login`, NewUserLoginHandler(store.User, store.ActiveUser, logger))

	// Private
	mux.HandleFunc(`GET /api/user/orders`, AuthMiddleware(NewUserOrdersHanlder(store.Order, logger)))
	mux.HandleFunc(`GET /api/user/balance`, AuthMiddleware(NewUserBalanceHanlder(store.User, logger)))
	mux.HandleFunc(`POST /api/user/orders`, AuthMiddleware(NewUploadUserOrdersHanlder(store.Order, store.OrdersQueue, logger)))
	mux.HandleFunc(`GET /api/user/withdrawals`, AuthMiddleware(NewUserWithdrawalsHanlder(store.Billing, logger)))
	mux.HandleFunc(`POST /api/user/balance/withdraw`, AuthMiddleware(NewUserBalanceWithdrawHanlder(store.User, store.Billing, logger)))

	// Manual debug
	mux.HandleFunc(`POST /api/user/orders/accrual`, AuthMiddleware(NewUploadUserOrdersAccrualHanlder(store.Order, logger)))

	return mux
}
