package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
)

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	slog.Info("Deleting user",
		slog.String("user_id", id.String()),
	)

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			slog.Warn("User not found",
				slog.String("user_id", id.String()),
			)
			return repository.ErrUserNotFound
		}
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return fmt.Errorf("%w: %w", repository.ErrFailedToGetUser, err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return fmt.Errorf("%w: %w", repository.ErrFailedToDeleteUser, err)
	}

	slog.Info("User deleted successfully",
		slog.String("user_id", id.String()),
	)

	return nil
}
