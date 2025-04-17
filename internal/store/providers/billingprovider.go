package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type BillingRepositoryProvider interface {
	Withdrawals(context.Context, int) ([]model.Withdraw, error)
	BalanceWithdraw(context.Context, *model.Withdraw, *model.User) error
}
