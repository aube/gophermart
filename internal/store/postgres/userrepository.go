package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aube/gophermart/internal/httperrors"
	"github.com/aube/gophermart/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepository ...
type UserRepository struct {
	db *sql.DB
}

// Register ...
func (r *UserRepository) Register(ctx context.Context, u *model.User) error {
	if err := u.Validate(); err != nil {
		return httperrors.NewValidationError(err)
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (login, encrypted_password) VALUES ($1, $2) ON CONFLICT (login) DO NOTHING RETURNING id",
		u.Login,
		u.EncryptedPassword,
	).Scan(&u.ID)

	// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return httperrors.NewConflictError()
	}

	return nil
}

// Login ...
func (r *UserRepository) Login(ctx context.Context, u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(
		"SELECT id, encrypted_password FROM users WHERE login = $1",
		u.Login,
	).Scan(
		&u.ID,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, httperrors.NewRecordNotFound()
		}

		return nil, err
	}

	if !u.ComparePassword(u.Password) {
		return nil, httperrors.NewLoginFailed()
	}

	return u, nil
}

// Balance ...
func (r *UserRepository) Balance(ctx context.Context, u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(
		"SELECT id, encrypted_password FROM users WHERE login = $1",
		u.Login,
	).Scan(
		&u.ID,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, httperrors.NewRecordNotFound()
		}

		return nil, err
	}

	if !u.ComparePassword(u.Password) {
		return nil, httperrors.NewLoginFailed()
	}

	return u, nil
}
