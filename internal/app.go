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

	client.NewServicePolling(store, config.AccrualSystemAddress)

	mux := api.NewRouter(store)

	return http.ListenAndServe(config.ServerAddress, mux)
}
