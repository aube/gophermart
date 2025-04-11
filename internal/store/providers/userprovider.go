package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type UserRepositoryProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, string) (*model.User, error)
	Orders(context.Context, string) (*model.User, error)
	Balance(context.Context, string) (*model.User, error)
	UploadOrders(context.Context, string) (*model.User, error)
	BalanceWithdraw(context.Context, string) (*model.User, error)
	Withdrawals(context.Context, string) (*model.User, error)
}
