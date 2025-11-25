package postgres

import (
	"context"
	"errors"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, phone_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, user.FullName, user.Email, user.PasswordHash, user.PhoneNumber).
		Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, phone_number, created_at
		FROM users
		WHERE email = $1
	`
	user := &domain.User{}
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.PhoneNumber, &user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, phone_number, created_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.PhoneNumber, &user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
