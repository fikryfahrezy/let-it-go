package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	setupTest(t)

	// Create test users
	users := []repository.User{
		{
			Name:     "User 1",
			Email:    "user1@example.com",
			Password: "password1",
		},
		{
			Name:     "User 2",
			Email:    "user2@example.com",
			Password: "password2",
		},
		{
			Name:     "User 3",
			Email:    "user3@example.com",
			Password: "password3",
		},
	}

	for _, user := range users {
		err := testRepository.Create(context.Background(), user)
		require.NoError(t, err)
	}

	// Test pagination
	result, err := testRepository.List(context.Background(), 2, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = testRepository.List(context.Background(), 2, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = testRepository.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestListEmpty(t *testing.T) {
	setupTest(t)

	result, err := testRepository.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
}
