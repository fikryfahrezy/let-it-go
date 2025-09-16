package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	slog.Info("Creating new user",
		slog.String("email", req.Email),
		slog.String("name", req.Name),
	)

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("Failed to check existing user",
			slog.String("error", err.Error()),
			slog.String("email", req.Email),
		)
		return CreateUserResponse{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		slog.Warn("User already exists",
			slog.String("email", req.Email),
		)
		return CreateUserResponse{}, fmt.Errorf("user with email %s already exists", req.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password",
			slog.String("error", err.Error()),
		)
		return CreateUserResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &repository.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		slog.Error("Failed to create user",
			slog.String("error", err.Error()),
			slog.String("email", req.Email),
		)
		return CreateUserResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	response := ToCreateUserResponse(user)
	slog.Info("User created successfully",
		slog.Int("user_id", user.ID),
		slog.String("email", user.Email),
	)

	return response, nil
}