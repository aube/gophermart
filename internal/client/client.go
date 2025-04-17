package client

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aube/gophermart/internal/store"
)

// NewServicePolling ...
func NewServicePolling(store store.Store, accSystemAddress string) error {
	// Create iterator for send orders to loyalty program

	go func() {
		for range time.Tick(108 * time.Millisecond) {
			sendOrderToService(store, accSystemAddress)
		}
	}()

	// Receive new orders and fill queue
	orders, err := store.Order.GetNewOrdersID()
	if err != nil {
		return err
	}

	for _, id := range orders {
		store.OrdersQueue.Enqueue(id)
	}

	return nil
}

func sendOrderToService(store store.Store, accSystemAddress string) {
	if store.OrdersQueue.IsEmpty() {
		return
	}

	id, err := store.OrdersQueue.Dequeue()
	if err != nil {
		return
	}

	fmt.Println("Dequeue order", id)

	err = store.Order.SetStatus(id, "PROCESSING")
	if err != nil {
		store.OrdersQueue.Enqueue(id)
		return
	}

	oa, err := request(accSystemAddress + "/api/orders/" + strconv.Itoa(id))

	fmt.Println("oa", oa)

	if errors.Is(err, errors.New("new")) {
		store.Order.SetStatus(id, "NEW")
		return
	}

	if errors.Is(err, errors.New("invalid")) {
		store.Order.SetStatus(id, "INVALID")
		return
	}

	if oa.Status == "INVALID" {
		store.Order.SetStatus(id, "INVALID")
		return
	}

	store.Order.SetAccrual(id, oa.Accrual)
}
