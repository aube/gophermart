package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type BillingRepositoryProvider interface {
	Balance(context.Context, *model.User) (*model.User, error)
	BalanceWithdraw(context.Context, *model.User, int) (*model.User, error)
	Withdrawals(context.Context, *model.User) (*model.User, error)
}
