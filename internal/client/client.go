package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aube/gophermart/internal/model"
)

type OrderProvider interface {
	Orders(context.Context, int) ([]model.Order, error)
	UploadOrders(context.Context, *model.Order) error
	GetNewOrdersID() ([]int, error)
	SetStatus(int, string) error
	SetAccrual(int, int) error
}

type OrdersQueueProvider interface {
	Enqueue(item int)
	Dequeue() (int, error)
	IsEmpty() bool
	Size() int
}

// NewServicePolling ...
func NewServicePolling(
	storeOrder OrderProvider,
	storeOrdersQueue OrdersQueueProvider,
	accSystemAddress string,
) error {
	// Create iterator for send orders to loyalty program

	go func() {
		for range time.Tick(108 * time.Millisecond) {
			sendOrderToService(storeOrder, storeOrdersQueue, accSystemAddress)
		}
	}()

	// Receive new orders and fill queue
	orders, err := storeOrder.GetNewOrdersID()
	if err != nil {
		return err
	}

	for _, id := range orders {
		storeOrdersQueue.Enqueue(id)
	}

	return nil
}

func sendOrderToService(
	storeOrder OrderProvider,
	storeOrdersQueue OrdersQueueProvider,
	accSystemAddress string,
) {
	if storeOrdersQueue.IsEmpty() {
		return
	}

	id, err := storeOrdersQueue.Dequeue()
	if err != nil {
		return
	}

	fmt.Println("Dequeue order", id)

	err = storeOrder.SetStatus(id, "PROCESSING")
	if err != nil {
		storeOrdersQueue.Enqueue(id)
		return
	}

	oa, err := request(accSystemAddress + "/api/orders/" + strconv.Itoa(id))

	fmt.Println("oa", oa)

	if errors.Is(err, errors.New("new")) {
		storeOrder.SetStatus(id, "NEW")
		return
	}

	if errors.Is(err, errors.New("invalid")) {
		storeOrder.SetStatus(id, "INVALID")
		return
	}

	if oa.Status == "INVALID" {
		storeOrder.SetStatus(id, "INVALID")
		return
	}

	storeOrder.SetAccrual(id, oa.Accrual)
}
