package memory

import (
	"github.com/aube/gophermart/internal/model"
	"github.com/aube/gophermart/internal/store/providers"
)

// Store ...
type MemoryStore struct {
	mem map[string]*model.User
	oq  []int
}

// ActiveUser ...
func (s *MemoryStore) ActiveUser() providers.ActiveUserRepositoryProvider {
	return &ActiveUserRepository{
		mem: s.mem,
	}
}

// OrdersQueue ...
func (s *MemoryStore) OrdersQueue() providers.OrdersQueueRepositoryProvider {
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
