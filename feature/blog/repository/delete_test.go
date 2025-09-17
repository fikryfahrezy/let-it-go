package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestDelete() {
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

	err = suite.repository.Delete(suite.ctx, createdBlog.ID)
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repository.GetByID(suite.ctx, createdBlog.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrBlogNotFound, err)
}

func (suite *BlogRepositoryTestSuite) TestDeleteNotFound() {
	randomID := uuid.New()
	err := suite.repository.Delete(suite.ctx, randomID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrBlogNotFound, err)
}

func TestBlogDeleteTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}