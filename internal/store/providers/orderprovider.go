package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type OrderRepositoryProvider interface {
	Orders(context.Context, int) ([]model.Order, error)
	UploadOrders(context.Context, *model.Order) error
	GetNewOrdersID(context.Context) ([]int, error)
	SetStatus(context.Context, int, string) error
	SetAccrual(context.Context, int, int) error
}
