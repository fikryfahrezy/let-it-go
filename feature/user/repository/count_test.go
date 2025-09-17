package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *UserRepositoryTestSuite) TestCount() {
	// Initially empty
	count, err := suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)

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
	}

	for _, user := range users {
		err := suite.repository.Create(suite.ctx, user)
		require.NoError(suite.T(), err)
	}

	count, err = suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)
}

func TestUserCountTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(UserRepositoryTestSuite))
}