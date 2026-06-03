package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakabrajadenta/go-explore-api/internal/model"
)

var ErrNotFound = errors.New("record not found")

type UserRepository interface {
	FindAll(ctx context.Context, page, perPage int) ([]model.User, int, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user model.User) (*model.User, error)
	Update(ctx context.Context, user model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, page, perPage int) ([]model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx,
		`SELECT id, username, email, full_name, phone, is_active, created_at, updated_at
		 FROM users ORDER BY id LIMIT $1 OFFSET $2`,
		perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
			&u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.scanOne(ctx,
		`SELECT id, username, email, full_name, phone, is_active, created_at, updated_at
		 FROM users WHERE id = $1`, id)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.scanOne(ctx,
		`SELECT id, username, email, full_name, phone, is_active, created_at, updated_at
		 FROM users WHERE username = $1`, username)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.scanOne(ctx,
		`SELECT id, username, email, full_name, phone, is_active, created_at, updated_at
		 FROM users WHERE email = $1`, email)
}

func (r *userRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (username, email, full_name, phone, is_active)
		 VALUES ($1, $2, $3, $4, TRUE)
		 RETURNING id, username, email, full_name, phone, is_active, created_at, updated_at`,
		user.Username, user.Email, user.FullName, user.Phone,
	).Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
		&u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}

func (r *userRepository) Update(ctx context.Context, user model.User) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`UPDATE users
		 SET username=$1, email=$2, full_name=$3, phone=$4, is_active=$5
		 WHERE id=$6
		 RETURNING id, username, email, full_name, phone, is_active, created_at, updated_at`,
		user.Username, user.Email, user.FullName, user.Phone, user.IsActive, user.ID,
	).Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
		&u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) scanOne(ctx context.Context, query string, args ...any) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
			&u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}
