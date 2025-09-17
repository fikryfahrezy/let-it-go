package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
)

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (UpdateUserResponse, error) {
	slog.Info("Updating user",
		slog.String("user_id", id.String()),
	)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			slog.Warn("User not found",
				slog.String("user_id", id.String()),
			)
			return UpdateUserResponse{}, repository.ErrUserNotFound
		}
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return UpdateUserResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToGetUser, err)
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			// If it's not a "not found" error, it's a real database issue
			if !errors.Is(err, repository.ErrUserNotFound) {
				slog.Error("Failed to check existing user",
					slog.String("error", err.Error()),
					slog.String("email", req.Email),
				)
				return UpdateUserResponse{}, fmt.Errorf("%w: %w", ErrFailedToCheckExistingUser, err)
			}
			// User not found is expected, continue
		} else if existingUser.ID != uuid.Nil && existingUser.ID != id {
			// User found with different ID, return business logic error
			slog.Warn("Email already exists",
				slog.String("email", req.Email),
			)
			return UpdateUserResponse{}, ErrUserAlreadyExists
		}

		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return UpdateUserResponse{}, fmt.Errorf("%w: %w", repository.ErrFailedToUpdateUser, err)
	}

	response := ToUpdateUserResponse(user)
	slog.Info("User updated successfully",
		slog.String("user_id", id.String()),
	)

	return response, nil
}
