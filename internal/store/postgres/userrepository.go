package postgres

import (
	"context"
	"database/sql"

	"github.com/aube/gophermart/internal/model"
)

// UserRepository ...
type UserRepository struct {
	db *sql.DB
}

// Register ...
func (r *UserRepository) Register(ctx context.Context, u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// Login ...
func (r *UserRepository) Login(ctx context.Context, u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(
		"SELECT id, email FROM users WHERE email = $1 and encrypted_password = $2",
		u.Email,
		u.EncryptedPassword,
	).Scan(
		&u.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// Orders ...
func (r *UserRepository) Orders(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT * FROM orders WHERE user_id = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// Orders ...
func (r *UserRepository) Balance(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT * FROM orders WHERE user_id = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// UploadOrders ...
func (r *UserRepository) UploadOrders(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT * FROM orders WHERE user_id = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// BalanceWithdraw ...
func (r *UserRepository) BalanceWithdraw(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT * FROM orders WHERE user_id = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// Withdrawals ...
func (r *UserRepository) Withdrawals(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT * FROM orders WHERE user_id = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
