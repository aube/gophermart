package app

import (
	"net/http"

	"github.com/aube/gophermart/internal/api"
	"github.com/aube/gophermart/internal/client"
	"github.com/aube/gophermart/internal/store"
)

// Start ...
func Start() error {
	config := NewConfig()

	store, err := store.NewStore(config.DatabaseDSN)

	if err != nil {
		panic(err)
	}

	client.NewServicePolling(
		store.Order,
		store.OrdersQueue,
		config.AccrualSystemAddress,
	)

	mux := api.NewRouter(
		store.ActiveUser,
		store.Billing,
		store.Order,
		store.OrdersQueue,
		store.User,
	)

	return http.ListenAndServe(config.ServerAddress, mux)
}
