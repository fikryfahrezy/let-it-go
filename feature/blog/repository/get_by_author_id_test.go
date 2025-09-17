package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestGetByAuthorID(t *testing.T) {
	authorID := setupTest(t)

	// Create another author
	author2ID := uuid.New()
	_, err := db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		author2ID, "Test Author 2", "author2@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}

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
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test get by author 1
	result, err := testRepository.GetByAuthorID(context.Background(), authorID, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 blogs for author 1, got %d", len(result))
	}

	// Test get by author 2
	result, err = testRepository.GetByAuthorID(context.Background(), author2ID, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 blog for author 2, got %d", len(result))
	}

	// Test pagination
	result, err = testRepository.GetByAuthorID(context.Background(), authorID, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 blog with pagination, got %d", len(result))
	}
}
