package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type OrderRepositoryProvider interface {
	Orders(context.Context, int) ([]model.Order, error)
	UploadOrders(context.Context, *model.Order) error
}
