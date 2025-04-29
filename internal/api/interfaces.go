package api

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type UserProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, *model.User) (*model.User, error)
	Balance(context.Context, *model.User) error
}

type ActiveUserProvider interface {
	Set(context.Context, *model.User) error
	Get(context.Context, string) (*model.User, bool)
}

type BillingProvider interface {
	Withdrawals(context.Context, int) ([]model.Withdraw, error)
	BalanceWithdraw(context.Context, *model.Withdraw, *model.User) error
}

type OrderProvider interface {
	Orders(context.Context, int) ([]model.Order, error)
	UploadOrders(context.Context, *model.Order) error
	GetNewOrdersID() ([]int, error)
	SetStatus(int, string) error
	SetAccrual(int, int) error
}

type AccrualService interface {
	AddWork(int)
}
