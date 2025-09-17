package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByID(t *testing.T) {
	setupTest(t)

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := testRepository.Create(context.Background(), user)
	require.NoError(t, err)

	// Get the created user by email to get the actual ID
	createdUser, err := testRepository.GetByEmail(context.Background(), user.Email)
	require.NoError(t, err)

	result, err := testRepository.GetByID(context.Background(), createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
}

func TestGetByIDNotFound(t *testing.T) {
	setupTest(t)

	randomID := uuid.New()
	_, err := testRepository.GetByID(context.Background(), randomID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
}
