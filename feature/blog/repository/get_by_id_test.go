package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestGetByID() {
	blog := Blog{
		Title:    "Test Blog",
		Content:  "This is a test blog content",
		AuthorID: suite.authorID,
		Status:   StatusDraft,
	}

	err := suite.repository.Create(suite.ctx, blog)
	require.NoError(suite.T(), err)

	createdBlog, err := suite.getBlogByTitle(blog.Title)
	require.NoError(suite.T(), err)

	result, err := suite.repository.GetByID(suite.ctx, createdBlog.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdBlog.ID, result.ID)
	assert.Equal(suite.T(), blog.Title, result.Title)
	assert.Equal(suite.T(), blog.Content, result.Content)
}

func (suite *BlogRepositoryTestSuite) TestGetByIDNotFound() {
	randomID := uuid.New()
	_, err := suite.repository.GetByID(suite.ctx, randomID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrBlogNotFound, err)
}

func TestBlogGetByIDTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}