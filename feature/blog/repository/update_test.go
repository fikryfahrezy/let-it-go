package repository

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestUpdate() {
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

	// Update blog
	createdBlog.Title = "Updated Blog Title"
	createdBlog.Content = "Updated content"
	createdBlog.Status = StatusPublished
	publishedAt := time.Now().UTC()
	createdBlog.PublishedAt = &publishedAt

	err = suite.repository.Update(suite.ctx, createdBlog)
	assert.NoError(suite.T(), err)

	// Verify update
	result, err := suite.repository.GetByID(suite.ctx, createdBlog.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Blog Title", result.Title)
	assert.Equal(suite.T(), "Updated content", result.Content)
	assert.Equal(suite.T(), StatusPublished, result.Status)
	assert.NotNil(suite.T(), result.PublishedAt)
}

func (suite *BlogRepositoryTestSuite) TestUpdateNotFound() {
	blog := Blog{
		ID:       uuid.New(),
		Title:    "Test Blog",
		Content:  "This is a test blog content",
		AuthorID: suite.authorID,
		Status:   StatusDraft,
	}

	err := suite.repository.Update(suite.ctx, blog)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrBlogNotFound, err)
}

func TestBlogUpdateTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}