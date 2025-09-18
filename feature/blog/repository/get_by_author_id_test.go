package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestGetByAuthorID(t *testing.T) {
	authorID := setupTest(t)

	// Create another author
	author2ID := uuid.New()
	_, err := db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		author2ID, "Test Author 2", "author2@example.com", "password")
	assert.NoError(t, err)

	// Create test blogs
	blogs := []repository.Blog{
		{
			Title:    "Author 1 Blog 1",
			Content:  "Content 1",
			AuthorID: authorID,
			Status:   repository.StatusDraft,
		},
		{
			Title:    "Author 1 Blog 2",
			Content:  "Content 2",
			AuthorID: authorID,
			Status:   repository.StatusPublished,
		},
		{
			Title:    "Author 2 Blog 1",
			Content:  "Content 3",
			AuthorID: author2ID,
			Status:   repository.StatusDraft,
		},
	}

	for _, blog := range blogs {
		err := testRepository.Create(context.Background(), blog)
		assert.NoError(t, err)
	}

	// Test get by author 1
	result, err := testRepository.GetByAuthorID(context.Background(), authorID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Test get by author 2
	result, err = testRepository.GetByAuthorID(context.Background(), author2ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// Test pagination
	result, err = testRepository.GetByAuthorID(context.Background(), authorID, 1, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
