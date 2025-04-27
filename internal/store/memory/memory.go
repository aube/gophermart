package memory

import (
	"context"

	"github.com/aube/gophermart/internal/model"
)

// ActiveUserProvider ...
type ActiveUserProvider interface {
	Set(context.Context, *model.User) error
	Get(context.Context, string) (*model.User, bool)
}

// OrdersQueueProvider ...
type OrdersQueueProvider interface {
	Enqueue(item int)
	Dequeue() (int, error)
	IsEmpty() bool
	Size() int
}

// Store ...
type MemoryStore struct {
	mem map[string]*model.User
	oq  []int
}

// ActiveUser ...
func (s *MemoryStore) ActiveUser() ActiveUserProvider {
	return &ActiveUserRepository{
		mem: s.mem,
	}
}

// OrdersQueue ...
func (s *MemoryStore) OrdersQueue() OrdersQueueProvider {
	return &OrdersQueueRepository{
		oq: s.oq,
	}
}

// New ...
func New() *MemoryStore {
	return &MemoryStore{
		mem: make(map[string]*model.User),
		oq:  []int{},
	}
}
