package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestCount() {
	// Initially empty
	count, err := suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)

	// Create test blogs
	blogs := []Blog{
		{
			Title:    "Blog 1",
			Content:  "Content 1",
			AuthorID: suite.authorID,
			Status:   StatusDraft,
		},
		{
			Title:    "Blog 2",
			Content:  "Content 2",
			AuthorID: suite.authorID,
			Status:   StatusPublished,
		},
	}

	for _, blog := range blogs {
		err := suite.repository.Create(suite.ctx, blog)
		require.NoError(suite.T(), err)
	}

	count, err = suite.repository.Count(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)
}

func (suite *BlogRepositoryTestSuite) TestCountByStatus() {
	// Create test blogs with different statuses
	blogs := []Blog{
		{
			Title:    "Draft Blog 1",
			Content:  "Content 1",
			AuthorID: suite.authorID,
			Status:   StatusDraft,
		},
		{
			Title:    "Draft Blog 2",
			Content:  "Content 2",
			AuthorID: suite.authorID,
			Status:   StatusDraft,
		},
		{
			Title:    "Published Blog",
			Content:  "Content 3",
			AuthorID: suite.authorID,
			Status:   StatusPublished,
		},
	}

	for _, blog := range blogs {
		err := suite.repository.Create(suite.ctx, blog)
		require.NoError(suite.T(), err)
	}

	// Test count by status
	count, err := suite.repository.CountByStatus(suite.ctx, StatusDraft)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)

	count, err = suite.repository.CountByStatus(suite.ctx, StatusPublished)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	count, err = suite.repository.CountByStatus(suite.ctx, StatusArchived)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func TestBlogCountTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}