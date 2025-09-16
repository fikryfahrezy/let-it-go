package service

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *userService) UpdateUser(ctx context.Context, id int, req UpdateUserRequest) (UpdateUserResponse, error) {
	slog.Info("Updating user",
		slog.Int("user_id", id),
	)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user by ID",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return UpdateUserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		slog.Warn("User not found",
			slog.Int("user_id", id),
		)
		return UpdateUserResponse{}, fmt.Errorf("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			slog.Error("Failed to check existing user",
				slog.String("error", err.Error()),
				slog.String("email", req.Email),
			)
			return UpdateUserResponse{}, fmt.Errorf("failed to check existing user: %w", err)
		}

		if existingUser != nil && existingUser.ID != id {
			slog.Warn("Email already exists",
				slog.String("email", req.Email),
			)
			return UpdateUserResponse{}, fmt.Errorf("user with email %s already exists", req.Email)
		}

		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return UpdateUserResponse{}, fmt.Errorf("failed to update user: %w", err)
	}

	response := ToUpdateUserResponse(user)
	slog.Info("User updated successfully",
		slog.Int("user_id", id),
	)

	return response, nil
}
