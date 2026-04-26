package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BogdanBratsky/eigth/internal/model"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (login, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err := r.db.QueryRowContext(ctx, query, user.Login, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, login, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, email).
		Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.PasswordHash,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserPostgres) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `
		SELECT id, login, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE login = $1;
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, login).
		Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.PasswordHash,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
