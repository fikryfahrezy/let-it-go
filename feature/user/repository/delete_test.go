package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *UserRepositoryTestSuite) TestDelete() {
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

	err = suite.repository.Delete(suite.ctx, createdUser.ID)
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repository.GetByID(suite.ctx, createdUser.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUserNotFound, err)
}

func (suite *UserRepositoryTestSuite) TestDeleteNotFound() {
	randomID := uuid.New()
	err := suite.repository.Delete(suite.ctx, randomID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUserNotFound, err)
}

func TestUserDeleteTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(UserRepositoryTestSuite))
}