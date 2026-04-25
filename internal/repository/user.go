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
	var user model.User
	query := `
		SELECT id, login, email, password_hash, is_active, created_at, updated_at
		FROM users 
		WHERE email = $1;
	`
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
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgres) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, login, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, id).
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
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgres) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	query := `
		SELECT id, login, created_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		u := model.User{}
		err := rows.Scan(
			&u.ID,
			&u.Login,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
