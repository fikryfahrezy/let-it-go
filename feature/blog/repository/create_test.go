package repository

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestCreate() {
	blog := Blog{
		Title:    "Test Blog",
		Content:  "This is a test blog content",
		AuthorID: suite.authorID,
		Status:   StatusDraft,
	}

	err := suite.repository.Create(suite.ctx, blog)
	assert.NoError(suite.T(), err)

	// Verify blog was created by getting all blogs and checking the title
	blogs, err := suite.repository.List(suite.ctx, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), blogs, 1)
	
	result := blogs[0]
	assert.NotEqual(suite.T(), uuid.Nil, result.ID)
	assert.Equal(suite.T(), blog.Title, result.Title)
	assert.Equal(suite.T(), blog.Content, result.Content)
	assert.Equal(suite.T(), blog.AuthorID, result.AuthorID)
	assert.Equal(suite.T(), blog.Status, result.Status)
	assert.NotZero(suite.T(), result.CreatedAt)
	assert.NotZero(suite.T(), result.UpdatedAt)
}

func (suite *BlogRepositoryTestSuite) TestCreateWithPublishedStatus() {
	publishedAt := time.Now().UTC()
	blog := Blog{
		Title:       "Published Blog",
		Content:     "This is a published blog content",
		AuthorID:    suite.authorID,
		Status:      StatusPublished,
		PublishedAt: &publishedAt,
	}

	err := suite.repository.Create(suite.ctx, blog)
	assert.NoError(suite.T(), err)

	result, err := suite.getBlogByTitle(blog.Title)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), StatusPublished, result.Status)
	assert.NotNil(suite.T(), result.PublishedAt)
}

func TestBlogCreateTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}