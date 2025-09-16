package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/fikryfahrezy/let-it-go/pkg/database"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}


func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, now, now)
	if err != nil {
		slog.Error("Failed to create user",
			slog.String("error", err.Error()),
			slog.String("email", user.Email),
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get last insert id",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = int(id)

	slog.Info("User created successfully",
		slog.Int("user_id", user.ID),
		slog.String("email", user.Email),
	)

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("Failed to get user by email",
			slog.String("error", err.Error()),
			slog.String("email", email),
		)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, now, user.ID)
	if err != nil {
		slog.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.Int("user_id", user.ID),
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	slog.Info("User updated successfully",
		slog.Int("user_id", user.ID),
	)

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete user",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	slog.Info("User deleted successfully",
		slog.Int("user_id", id),
	)

	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Failed to list users",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan user row",
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating user rows",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		slog.Error("Failed to count users",
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}