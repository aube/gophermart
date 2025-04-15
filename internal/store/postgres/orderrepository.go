package postgres

import (
	"context"
	"database/sql"
	"errors"
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

// GetNewOrdersID ...
func (r *OrderRepository) GetNewOrdersID(ctx context.Context) ([]int, error) {

	rows, err := r.db.QueryContext(
		ctx,
		"select id from orders where status='NEW'",
	)
	if err != nil {
		return nil, errors.New("Orders select error")
	}
	defer rows.Close()

	var result []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return nil, errors.New("Order Scan error")
		}

		result = append(result, id)
	}

	return result, nil
}

// SetStatus ...
func (r *OrderRepository) SetStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(
		ctx,
		"update orders set status=$1 where id=$2",
		status,
		id,
	)

	if err != nil {
		return errors.New("Order change status error")
	}

	return nil
}

// SetAccrual ...
func (r *OrderRepository) SetAccrual(ctx context.Context, id int, accrual int) error {

	var userID int
	var balance int

	// Select user data by order_id
	if err := r.db.QueryRow(
		`SELECT o.user_id as userID, u.balance as balance
		FROM orders as o
		LEFT JOIN users as u ON o.user_id = u.id
		WHERE o.id = $1`,
		id,
	).Scan(
		&userID,
		&balance,
	); err != nil {
		return err
	}

	// TX begin
	tx, err := r.db.Begin()
	if err != nil {
		return httperrors.NewServerError(err)
	}

	// Change order status and accrual
	_, err = tx.ExecContext(
		ctx,
		"update orders set status=$1, accrual=$2 where id=$3",
		"PROCESSED",
		accrual,
		id,
	)

	if err != nil {
		tx.Rollback()
		return errors.New("Order change status error")
	}

	// Change user balance
	_, err = tx.ExecContext(
		ctx,
		"update users set balance=$1 where id=$3",
		balance+accrual,
		userID,
	)

	if err != nil {
		tx.Rollback()
		return errors.New("User change balance error")
	}

	tx.Commit()

	return nil
}
