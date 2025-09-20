package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/stretchr/testify/assert"
)

func TestCountByStatus(t *testing.T) {
	authorID := setupTest(t)

	// Create test blogs with different statuses
	blogs := []repository.Blog{
		{
			Title:    "Draft Blog 1",
			Content:  "Content 1",
			AuthorID: authorID,
			Status:   repository.StatusDraft,
		},
		{
			Title:    "Draft Blog 2",
			Content:  "Content 2",
			AuthorID: authorID,
			Status:   repository.StatusDraft,
		},
		{
			Title:    "Published Blog",
			Content:  "Content 3",
			AuthorID: authorID,
			Status:   repository.StatusPublished,
		},
	}

	for _, blog := range blogs {
		err := testRepository.Create(context.Background(), blog)
		assert.NoError(t, err)
	}

	// Test count by status
	count, err := testRepository.CountByStatus(context.Background(), repository.StatusDraft)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = testRepository.CountByStatus(context.Background(), repository.StatusPublished)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = testRepository.CountByStatus(context.Background(), repository.StatusArchived)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
