package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/aube/gophermart/internal/store/memory"
	"github.com/aube/gophermart/internal/store/postgres"
	"github.com/aube/gophermart/internal/store/providers"
)

type Store struct {
	User        providers.UserRepositoryProvider
	ActiveUser  providers.ActiveUserRepositoryProvider
	Order       providers.OrderRepositoryProvider
	Billing     providers.BillingRepositoryProvider
	OrdersQueue providers.OrdersQueueRepositoryProvider
}

func NewStore(config string) (Store, error) {
	db, err := NewDB(config)
	if err != nil {
		return Store{}, err
	}

	store := Store{
		User:        postgres.New(db).User(),
		Order:       postgres.New(db).Order(),
		Billing:     postgres.New(db).Billing(),
		ActiveUser:  memory.New().ActiveUser(),
		OrdersQueue: memory.New().OrdersQueue(),
	}

	return store, nil
}

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("PostgreSQL database connection established")

	runPostgresMigrations(db)

	return db, nil
}
