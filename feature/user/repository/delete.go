package repository

import (
	"context"
	"fmt"
	"log/slog"
)

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