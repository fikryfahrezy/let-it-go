package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	setupTest(t)

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := testRepository.Create(context.Background(), user)
	require.NoError(t, err)

	// Get the created user to get the actual ID
	createdUser, err := testRepository.GetByEmail(context.Background(), user.Email)
	require.NoError(t, err)

	// Update user
	createdUser.Name = "John Updated"
	createdUser.Email = "john.updated@example.com"
	err = testRepository.Update(context.Background(), createdUser)
	assert.NoError(t, err)

	// Verify update
	result, err := testRepository.GetByID(context.Background(), createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "John Updated", result.Name)
	assert.Equal(t, "john.updated@example.com", result.Email)
}

func TestUpdateNotFound(t *testing.T) {
	setupTest(t)

	user := repository.User{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := testRepository.Update(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
}
