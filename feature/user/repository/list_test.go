package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *UserRepositoryTestSuite) TestList() {
	// Create test users
	users := []User{
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
		err := suite.repository.Create(suite.ctx, user)
		require.NoError(suite.T(), err)
	}

	// Test pagination
	result, err := suite.repository.List(suite.ctx, 2, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	result, err = suite.repository.List(suite.ctx, 2, 1)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	result, err = suite.repository.List(suite.ctx, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)
}

func (suite *UserRepositoryTestSuite) TestListEmpty() {
	result, err := suite.repository.List(suite.ctx, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

func TestUserListTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(UserRepositoryTestSuite))
}