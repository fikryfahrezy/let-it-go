package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *UserRepositoryTestSuite) TestGetByID() {
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := suite.repository.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// Get the created user by email to get the actual ID
	createdUser, err := suite.repository.GetByEmail(suite.ctx, user.Email)
	require.NoError(suite.T(), err)

	result, err := suite.repository.GetByID(suite.ctx, createdUser.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdUser.ID, result.ID)
	assert.Equal(suite.T(), user.Name, result.Name)
	assert.Equal(suite.T(), user.Email, result.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByIDNotFound() {
	randomID := uuid.New()
	_, err := suite.repository.GetByID(suite.ctx, randomID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUserNotFound, err)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := suite.repository.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	result, err := suite.repository.GetByEmail(suite.ctx, user.Email)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, result.ID)
	assert.Equal(suite.T(), user.Name, result.Name)
	assert.Equal(suite.T(), user.Email, result.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByEmailNotFound() {
	_, err := suite.repository.GetByEmail(suite.ctx, "nonexistent@example.com")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUserNotFound, err)
}

func TestUserGetTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(UserRepositoryTestSuite))
}