package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	authorID := setupTest(t)

	// Initially empty
	count, err := testRepository.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create test blogs
	blogs := []repository.Blog{
		{
			Title:    "Blog 1",
			Content:  "Content 1",
			AuthorID: authorID,
			Status:   repository.StatusDraft,
		},
		{
			Title:    "Blog 2",
			Content:  "Content 2",
			AuthorID: authorID,
			Status:   repository.StatusPublished,
		},
	}

	for _, blog := range blogs {
		err := testRepository.Create(context.Background(), blog)
		assert.NoError(t, err)
	}

	count, err = testRepository.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
