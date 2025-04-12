package postgres

import (
	"context"
	"database/sql"

	"github.com/aube/gophermart/internal/model"
)

// OrderRepository ...
type OrderRepository struct {
	db *sql.DB
}

// Orders ...
func (r *OrderRepository) Orders(ctx context.Context, u *model.User) (*model.User, error) {
	if err := r.db.QueryRow(
		"SELECT * FROM orders WHERE user_id = $1",
		u.ID,
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

// UploadOrders ...
func (r *OrderRepository) UploadOrders(ctx context.Context, u *model.User, orderID int) (*model.User, error) {
	if err := r.db.QueryRow(
		"insert into orders set user_id = $1, order_id = $2",
		u.ID,
		orderID,
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
