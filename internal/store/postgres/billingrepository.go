package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

// BillingRepository ...
type BillingRepository struct {
	db *sql.DB
}

// BalanceWithdraw ...
func (r *BillingRepository) BalanceWithdraw(ctx context.Context, wd *model.Withdraw, u *model.User) error {
	if u.Balance < wd.Amount {
		return httperrors.NewNotEnoughMoneyError()
	}

	tx, err := r.db.Begin()
	if err != nil {
		return httperrors.NewServerError(err)
	}

	var newID int
	err = tx.QueryRowContext(
		ctx,
		"insert into withdrawals set user_id = $1, amount = $2 RETURNING id",
		u.ID,
		wd.Amount,
	).Scan(&newID)

	if err != nil {
		tx.Rollback()
		return httperrors.NewServerError(err)
	}

	if newID == 0 {
		tx.Rollback()
		return httperrors.NewServerError(errors.New("withdraw error"))
	}

	err = tx.QueryRowContext(
		ctx,
		"update users set balance = $1, withdrawn = $2 where id = $3",
		u.Balance-wd.Amount,
		u.Withdrawn+wd.Amount,
		u.ID,
	).Scan(&newID)

	if err != nil {
		tx.Rollback()
		return httperrors.NewServerError(err)
	}

	return tx.Commit()

}

// Withdrawals ...
func (r *BillingRepository) Withdrawals(ctx context.Context, userID int) ([]model.Withdraw, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"select order_id, amount/100 as sum, created_at from billing where user_id=$1",
		userID,
	)

	if err != nil {
		return []model.Withdraw{}, httperrors.NewServerError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return []model.Withdraw{}, httperrors.NewServerError(err)
	}

	var result []model.Withdraw

	for rows.Next() {
		var wd model.Withdraw
		err := rows.Scan(&wd.OrderID, &wd.Sum, &wd.ProcessedAt)

		if err != nil {
			return []model.Withdraw{}, httperrors.NewServerError(err)
		}

		result = append(result, wd)
	}

	return result, nil
}
