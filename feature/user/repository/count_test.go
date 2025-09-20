package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCount(t *testing.T) {
	setupTest(t)

	// Initially empty
	count, err := testRepository.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

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
	}

	for _, user := range users {
		err := testRepository.Create(context.Background(), user)
		require.NoError(t, err)
	}

	count, err = testRepository.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
