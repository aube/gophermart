package api

import (
	"net/http"

	"github.com/aube/gophermart/internal/logger"
)

func NewRouter(
	storeActiveUser ActiveUserProvider,
	storeBilling BillingProvider,
	storeOrder OrderProvider,
	storeOrdersQueue OrdersQueueProvider,
	storeUser UserProvider,
) *http.ServeMux {
	mux := http.NewServeMux()
	logger := logger.New()
	AuthMiddleware := NewAuthMiddleware(storeActiveUser, logger)

	// Public
	mux.HandleFunc(`POST /api/user/register`, NewUserRegisterHandler(storeUser, storeActiveUser, logger))
	mux.HandleFunc(`POST /api/user/login`, NewUserLoginHandler(storeUser, storeActiveUser, logger))

	// Private
	mux.HandleFunc(`GET /api/user/orders`, AuthMiddleware(NewUserOrdersHanlder(storeOrder, logger)))
	mux.HandleFunc(`GET /api/user/balance`, AuthMiddleware(NewUserBalanceHanlder(storeUser, logger)))
	mux.HandleFunc(`POST /api/user/orders`, AuthMiddleware(NewUploadUserOrdersHanlder(storeOrder, storeOrdersQueue, logger)))
	mux.HandleFunc(`GET /api/user/withdrawals`, AuthMiddleware(NewUserWithdrawalsHanlder(storeBilling, logger)))
	mux.HandleFunc(`POST /api/user/balance/withdraw`, AuthMiddleware(NewUserBalanceWithdrawHanlder(storeUser, storeBilling, logger)))

	// Manual debug
	mux.HandleFunc(`POST /api/user/orders/accrual`, AuthMiddleware(NewUploadUserOrdersAccrualHanlder(storeOrder, logger)))

	return mux
}
