package memory

import (
	"github.com/aube/gophermart/internal/model"
	"github.com/aube/gophermart/internal/store/providers"
)

// Store ...
type MemoryStore struct {
	mem map[string]*model.User
}

// ActiveUser ...
func (s *MemoryStore) ActiveUser() providers.ActiveUserRepositoryProvider {
	return &ActiveUserRepository{
		mem: s.mem,
	}
}

// New ...
func New() *MemoryStore {
	return &MemoryStore{
		mem: make(map[string]*model.User),
	}
}
