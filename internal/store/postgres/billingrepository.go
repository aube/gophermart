package postgres

import (
	"context"
	"database/sql"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
)

// BillingRepository ...
type BillingRepository struct {
	db *sql.DB
}

// Balance ...
func (r *BillingRepository) Balance(ctx context.Context, u *model.User) (*model.User, error) {
	if err := r.db.QueryRow(
		"SELECT * FROM users WHERE user_id = $1",
		u.ID,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, httperrors.NewRecordNotFound()
		}

		return nil, err
	}

	return u, nil
}

// BalanceWithdraw ...
func (r *BillingRepository) BalanceWithdraw(ctx context.Context, u *model.User, points int) (*model.User, error) {
	if err := r.db.QueryRow(
		"insert into withdrawals set user_id = $1, accrual = $2",
		u.ID,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, httperrors.NewRecordNotFound()
		}

		return nil, err
	}

	return u, nil
}

// Withdrawals ...
func (r *BillingRepository) Withdrawals(ctx context.Context, u *model.User) (*model.User, error) {
	if err := r.db.QueryRow(
		"SELECT * FROM withdrawals WHERE user_id = $1",
		u.ID,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, httperrors.NewRecordNotFound()
		}

		return nil, err
	}

	return u, nil
}
