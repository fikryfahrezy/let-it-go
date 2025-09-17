package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
)

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (GetUserResponse, error) {
	slog.Info("Getting user by ID",
		slog.String("user_id", id.String()),
	)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			slog.Warn("User not found",
				slog.String("user_id", id.String()),
			)
			return GetUserResponse{}, repository.ErrUserNotFound
		}
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return GetUserResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToGetUser, err)
	}

	response := ToGetUserResponse(user)
	return response, nil
}
