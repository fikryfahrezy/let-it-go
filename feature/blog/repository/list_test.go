package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestList(t *testing.T) {
	authorID := setupTest(t)

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
		{
			Title:    "Blog 3",
			Content:  "Content 3",
			AuthorID: authorID,
			Status:   repository.StatusDraft,
		},
	}

	for _, blog := range blogs {
		err := testRepository.Create(context.Background(), blog)
		assert.NoError(t, err)
	}

	// Test pagination
	result, err := testRepository.List(context.Background(), 2, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = testRepository.List(context.Background(), 2, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	result, err = testRepository.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
}
