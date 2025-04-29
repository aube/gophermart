package app

import (
	"net/http"

	"github.com/aube/gophermart/internal/api"
	"github.com/aube/gophermart/internal/client"
	"github.com/aube/gophermart/internal/store"
	"github.com/aube/gophermart/internal/workerpool"
)

// Start ...
func Start() error {
	config := NewConfig()

	store, err := store.NewStore(config.DatabaseDSN)

	if err != nil {
		panic(err)
	}

	accrualClient := client.New(config.AccrualSystemAddress, store.Order)
	dispatcher := workerpool.New(3, accrualClient.SendOrder)
	defer dispatcher.Close()

	// Receive new orders from database
	orders, err := store.Order.GetNewOrdersID()
	if err != nil {
		return err
	}

	for _, id := range orders {
		dispatcher.AddWork(id)
	}

	mux := api.NewRouter(
		store.ActiveUser,
		store.Billing,
		store.Order,
		store.User,
		dispatcher,
	)

	return http.ListenAndServe(config.ServerAddress, mux)
}
