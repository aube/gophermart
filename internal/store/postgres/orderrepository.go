package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

// OrderRepository ...
type OrderRepository struct {
	db *sql.DB
}

// Orders ...
func (r *OrderRepository) Orders(ctx context.Context, userID int) ([]model.Order, error) {

	rows, err := r.db.QueryContext(
		ctx,
		"select id, accrual, status, created_at from orders where user_id=$1",
		userID,
	)
	if err != nil {
		return []model.Order{}, httperrors.NewServerError(err)
	}
	defer rows.Close()

	var result []model.Order

	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.ID, &o.Accrual, &o.Status, &o.UploadedAT)

		if err != nil {
			return []model.Order{}, httperrors.NewServerError(err)
		}

		result = append(result, o)
	}

	return result, nil
}

// UploadOrders ...
func (r *OrderRepository) UploadOrders(ctx context.Context, o *model.Order) error {
	if err := o.CreateValidate(); err != nil {
		return httperrors.NewValidationError(err)
	}

	var newID int
	r.db.QueryRowContext(
		ctx,
		"INSERT INTO orders (id, user_id) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING RETURNING id",
		o.ID,
		o.UserID,
	).Scan(&newID)

	fmt.Println(newID)

	if newID == 0 {
		return httperrors.NewAlreadyUploadedError()
	}

	return nil
}
