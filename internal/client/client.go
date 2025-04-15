package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aube/gophermart/internal/store"
)

// NewServicePolling ...
func NewServicePolling(store store.Store, accSystemAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Receive new orders and fill queue
	orders, err := store.Order.GetNewOrdersID(ctx)
	if err != nil {
		return err
	}

	for id := range orders {
		store.OrdersQueue.Enqueue(id)
	}

	// Create iterator for send orders to loyalty program
	setInterval(1*time.Second, func() {
		sendOrderToService(store, accSystemAddress)
	})

	return nil
}

func sendOrderToService(store store.Store, accSystemAddress string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fmt.Println("Tick at", time.Now().Format("15:04:05"))
	if store.OrdersQueue.IsEmpty() {
		return
	}

	id, err := store.OrdersQueue.Dequeue()
	if err != nil {
		return
	}

	err = store.Order.SetStatus(ctx, id, "PROCESSING")
	if err != nil {
		store.OrdersQueue.Enqueue(id)
		return
	}

	oa, err := request(accSystemAddress + "/" + strconv.Itoa(id))

	if errors.Is(err, errors.New("new")) {
		store.Order.SetStatus(ctx, id, "NEW")
		return
	}

	if errors.Is(err, errors.New("invalid")) {
		store.Order.SetStatus(ctx, id, "INVALID")
		return
	}

	if oa.Status == "INVALID" {
		store.Order.SetStatus(ctx, id, "INVALID")
		return
	}

	store.Order.SetAccrual(ctx, id, oa.Accrual)
}

func setInterval(interval time.Duration, action func()) chan bool {
	ticker := time.NewTicker(interval)
	stopChan := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				action()
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	return stopChan
}
