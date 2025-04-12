package postgres

import (
	"database/sql"

	"github.com/aube/gophermart/internal/store/providers"
)

// Store ...
type SQLStore struct {
	db *sql.DB
}

// User ...
func (s *SQLStore) User() providers.UserRepositoryProvider {
	return &UserRepository{
		db: s.db,
	}
}

// Order ...
func (s *SQLStore) Order() providers.OrderRepositoryProvider {
	return &OrderRepository{
		db: s.db,
	}
}

// Billing ...
func (s *SQLStore) Billing() providers.BillingRepositoryProvider {
	return &BillingRepository{
		db: s.db,
	}
}

// New ...
func New(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}
