package providers

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

type UserRepositoryProvider interface {
	Register(context.Context, *model.User) error
	Login(context.Context, *model.User) (*model.User, error)
}
