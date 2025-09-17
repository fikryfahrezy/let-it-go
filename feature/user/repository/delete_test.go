package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := testRepository.Create(ctx, user)
	require.NoError(t, err)

	// Get the created user to get the actual ID
	createdUser, err := testRepository.GetByEmail(ctx, user.Email)
	require.NoError(t, err)

	err = testRepository.Delete(ctx, createdUser.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = testRepository.GetByID(ctx, createdUser.ID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
}

func TestDeleteNotFound(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	randomID := uuid.New()
	err := testRepository.Delete(ctx, randomID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
}
