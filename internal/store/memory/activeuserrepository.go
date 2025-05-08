package memory

import (
	"context"
	"sync"

	"github.com/aube/gophermart/internal/model"
)

// ActiveUserRepository ...
type ActiveUserRepository struct {
	mu  sync.Mutex
	mem map[string]*model.User
}

// Set ...
func (r *ActiveUserRepository) Set(ctx context.Context, u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.mem[u.RandomHash] = u

	return nil
}

// Get ...
func (r *ActiveUserRepository) Get(ctx context.Context, RandomHash string) (*model.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.mem[RandomHash]

	return user, ok
}
