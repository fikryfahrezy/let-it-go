package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByEmail(t *testing.T) {
	setupTest(t)

	user := repository.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := testRepository.Create(context.Background(), user)
	require.NoError(t, err)

	result, err := testRepository.GetByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
}

func TestGetByEmailNotFound(t *testing.T) {
	setupTest(t)

	_, err := testRepository.GetByEmail(context.Background(), "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
}