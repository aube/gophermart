package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type UserRepositoryProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, *model.User) (*model.User, error)
	Orders(context.Context, *model.User) (*model.User, error)
	Balance(context.Context, *model.User) (*model.User, error)
	UploadOrders(context.Context, *model.User, int) (*model.User, error)
	BalanceWithdraw(context.Context, *model.User, int) (*model.User, error)
	Withdrawals(context.Context, *model.User) (*model.User, error)
}
