package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *UserRepositoryTestSuite) TestCreate() {
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword",
	}

	err := suite.repository.Create(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Verify user was created by email since Create generates new ID
	result, err := suite.repository.GetByEmail(suite.ctx, user.Email)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, result.ID)
	assert.Equal(suite.T(), user.Name, result.Name)
	assert.Equal(suite.T(), user.Email, result.Email)
	assert.Equal(suite.T(), user.Password, result.Password)
	assert.NotZero(suite.T(), result.CreatedAt)
	assert.NotZero(suite.T(), result.UpdatedAt)
}

func (suite *UserRepositoryTestSuite) TestCreateDuplicateEmail() {
	user1 := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashedpassword1",
	}

	user2 := User{
		Name:     "Jane Doe",
		Email:    "john@example.com", // Same email
		Password: "hashedpassword2",
	}

	err := suite.repository.Create(suite.ctx, user1)
	assert.NoError(suite.T(), err)

	err = suite.repository.Create(suite.ctx, user2)
	assert.Error(suite.T(), err)
}

func TestUserCreateTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(UserRepositoryTestSuite))
}