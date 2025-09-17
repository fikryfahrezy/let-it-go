package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestGetByStatus() {
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
		{
			Title:    "Archived Blog",
			Content:  "Content 4",
			AuthorID: suite.authorID,
			Status:   StatusArchived,
		},
	}

	for _, blog := range blogs {
		err := suite.repository.Create(suite.ctx, blog)
		require.NoError(suite.T(), err)
	}

	// Test get blogs by status
	result, err := suite.repository.GetByStatus(suite.ctx, StatusDraft, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	result, err = suite.repository.GetByStatus(suite.ctx, StatusPublished, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)

	result, err = suite.repository.GetByStatus(suite.ctx, StatusArchived, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)

	// Test pagination
	result, err = suite.repository.GetByStatus(suite.ctx, StatusDraft, 1, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
}

func TestBlogGetByStatusTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}