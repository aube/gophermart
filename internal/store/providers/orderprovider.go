package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type OrderRepositoryProvider interface {
	Orders(context.Context, *model.User) (*model.User, error)
	UploadOrders(context.Context, *model.User, int) (*model.User, error)
}
