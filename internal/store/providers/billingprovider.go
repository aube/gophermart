package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type BillingRepositoryProvider interface {
	Balance(context.Context, *model.User) (*model.User, error)
	Withdrawals(ctx context.Context, u *model.User) ([]model.Withdraw, error)
	BalanceWithdraw(context.Context, *model.Withdraw, *model.User) error
}
