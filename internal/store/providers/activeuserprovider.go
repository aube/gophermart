package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type ActiveUserRepositoryProvider interface {
	Set(context.Context, *model.User) error
	Get(context.Context, string) (*model.User, bool)
}
