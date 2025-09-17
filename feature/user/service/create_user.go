package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	s.log.Info("Creating new user",
		slog.String("email", req.Email),
		slog.String("name", req.Name),
	)

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// If it's not a "not found" error, it's a real database issue
		if !errors.Is(err, repository.ErrUserNotFound) {
			return CreateUserResponse{}, err
		}
		// User not found is expected for creation, continue
	} else if existingUser.ID != uuid.Nil {
		// User found, return business logic error
		s.log.Warn("User already exists",
			slog.String("email", req.Email),
		)
		return CreateUserResponse{}, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed to hash password",
			slog.String("error", err.Error()),
		)
		return CreateUserResponse{}, fmt.Errorf("%w: %w", ErrFailedToHashPassword, err)
	}

	user := repository.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return CreateUserResponse{}, err
	}

	response := ToCreateUserResponse(user)
	s.log.Info("User created successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("email", user.Email),
	)

	return response, nil
}
