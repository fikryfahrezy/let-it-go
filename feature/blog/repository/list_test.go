package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestList() {
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
		{
			Title:    "Blog 3",
			Content:  "Content 3",
			AuthorID: suite.authorID,
			Status:   StatusArchived,
		},
	}

	for _, blog := range blogs {
		err := suite.repository.Create(suite.ctx, blog)
		require.NoError(suite.T(), err)
	}

	// Test pagination
	result, err := suite.repository.List(suite.ctx, 2, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	result, err = suite.repository.List(suite.ctx, 2, 1)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	result, err = suite.repository.List(suite.ctx, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)
}

func TestBlogListTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}