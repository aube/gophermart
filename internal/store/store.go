package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/aube/gophermart/internal/store/postgres"
	"github.com/aube/gophermart/internal/store/providers"
)

type Store struct {
	User providers.UserRepositoryProvider
}

func NewStore(config string) (Store, error) {
	db, err := NewDB(config)

	store := Store{
		User: postgres.New(db).User(),
	}

	return store, err
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
