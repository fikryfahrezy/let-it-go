package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestGetByStatus(t *testing.T) {
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
		{
			Title:    "Archived Blog",
			Content:  "Content 4",
			AuthorID: authorID,
			Status:   repository.StatusArchived,
		},
	}

	for _, blog := range blogs {
		err := testRepository.Create(context.Background(), blog)
		assert.NoError(t, err)
	}

	// Test get by draft status
	result, err := testRepository.GetByStatus(context.Background(), repository.StatusDraft, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Test get by published status
	result, err = testRepository.GetByStatus(context.Background(), repository.StatusPublished, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// Test get by archived status
	result, err = testRepository.GetByStatus(context.Background(), repository.StatusArchived, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// Test pagination
	result, err = testRepository.GetByStatus(context.Background(), repository.StatusDraft, 1, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
