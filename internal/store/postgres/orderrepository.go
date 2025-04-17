package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

	if err := rows.Err(); err != nil {
		return []model.Order{}, httperrors.NewServerError(err)
	}

	var result []model.Order

	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.ID, &o.Accrual, &o.Status, &o.UploadedAT)

		if err != nil {
			return []model.Order{}, httperrors.NewServerError(err)
		}

		o.Accrual /= 100
		result = append(result, o)
	}

	return result, nil
}

// UploadOrders ...
func (r *OrderRepository) UploadOrders(ctx context.Context, o *model.Order) error {
	if err := o.UploadOrderValidate(); err != nil {
		return httperrors.NewValidationError(err)
	}

	if !o.LuhnCheck(o.ID) {
		return httperrors.NewOrderNumberError()
	}

	var newID int
	r.db.QueryRowContext(
		ctx,
		"INSERT INTO orders (id, user_id) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING RETURNING id",
		o.ID,
		o.UserID,
	).Scan(&newID)

	if newID > 0 {
		return nil
	}

	var userID int
	r.db.QueryRowContext(
		ctx,
		"select user_id from orders where id=$1",
		o.ID,
	).Scan(&userID)

	if userID == o.UserID {
		return httperrors.NewAlreadyUploadedByMeError()
	}

	return httperrors.NewAlreadyUploadedAnotherError()
}

// GetNewOrdersID ...
func (r *OrderRepository) GetNewOrdersID() ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(
		ctx,
		"select id from orders where status='NEW'",
	)
	if err != nil {
		return nil, errors.New("orders select error")
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return []int{}, httperrors.NewServerError(err)
	}

	var result []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return nil, errors.New("order Scan error")
		}

		result = append(result, id)
	}

	return result, nil
}

// SetStatus ...
func (r *OrderRepository) SetStatus(id int, status string) error {

	fmt.Println("SetStatus", id, status)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		"update orders set status=$1 where id=$2",
		status,
		id,
	)

	if err != nil {
		return errors.New("order change status error")
	}

	return nil
}

// SetAccrual ...
func (r *OrderRepository) SetAccrual(id int, accrual int) error {

	fmt.Println("SetAccrual", id, accrual)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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
		return errors.New("urder change status error")
	}

	// Change user balance
	_, err = tx.ExecContext(
		ctx,
		"update users set balance=$1 where id=$2",
		balance+accrual,
		userID,
	)

	if err != nil {
		tx.Rollback()
		return errors.New("user change balance error")
	}

	tx.Commit()

	return nil
}
