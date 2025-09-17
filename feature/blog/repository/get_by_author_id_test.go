package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func (suite *BlogRepositoryTestSuite) TestGetByAuthorID() {
	// Create another author
	author2ID := uuid.New()
	_, err := suite.db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		author2ID, "Author 2", "author2@example.com", "password")
	require.NoError(suite.T(), err)

	// Create blogs for different authors
	blogsAuthor1 := []Blog{
		{
			Title:    "Author 1 Blog 1",
			Content:  "Content 1",
			AuthorID: suite.authorID,
			Status:   StatusDraft,
		},
		{
			Title:    "Author 1 Blog 2",
			Content:  "Content 2",
			AuthorID: suite.authorID,
			Status:   StatusPublished,
		},
	}

	blogsAuthor2 := []Blog{
		{
			Title:    "Author 2 Blog 1",
			Content:  "Content 3",
			AuthorID: author2ID,
			Status:   StatusDraft,
		},
	}

	for _, blog := range append(blogsAuthor1, blogsAuthor2...) {
		err := suite.repository.Create(suite.ctx, blog)
		require.NoError(suite.T(), err)
	}

	// Test get blogs by author
	result, err := suite.repository.GetByAuthorID(suite.ctx, suite.authorID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	result, err = suite.repository.GetByAuthorID(suite.ctx, author2ID, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)

	// Test pagination
	result, err = suite.repository.GetByAuthorID(suite.ctx, suite.authorID, 1, 0)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
}

func TestBlogGetByAuthorIDTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(BlogRepositoryTestSuite))
}