package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *UserRepositoryTestSuite) TestUpdate() {
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := suite.repository.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// Get the created user to get the actual ID
	createdUser, err := suite.repository.GetByEmail(suite.ctx, user.Email)
	require.NoError(suite.T(), err)

	// Update user
	createdUser.Name = "John Updated"
	createdUser.Email = "john.updated@example.com"
	err = suite.repository.Update(suite.ctx, createdUser)
	assert.NoError(suite.T(), err)

	// Verify update
	result, err := suite.repository.GetByID(suite.ctx, createdUser.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "John Updated", result.Name)
	assert.Equal(suite.T(), "john.updated@example.com", result.Email)
}

func (suite *UserRepositoryTestSuite) TestUpdateNotFound() {
	user := User{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := suite.repository.Update(suite.ctx, user)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUserNotFound, err)
}

func TestUserUpdateTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(UserRepositoryTestSuite))
}