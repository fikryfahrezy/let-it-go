package service

import "github.com/fikryfahrezy/let-it-go/pkg/app_error"

// Business logic errors (service-specific only)
var (
	// Business-specific errors (not in repository layer)
	ErrUserAlreadyExists = app_error.New("USER-USER_ALREADY_EXISTS", "user with email already exists")

	// Authentication errors
	ErrInvalidCredentials = app_error.New("USER-INVALID_CREDENTIALS", "invalid credentials")
	ErrFailedToHashPassword = app_error.New("USER-FAILED_TO_HASH_PASSWORD", "failed to hash password")

	// Validation errors (service-specific)
	ErrFailedToCheckExistingUser = app_error.New("USER-FAILED_TO_CHECK_EXISTING_USER", "failed to check existing user")
)