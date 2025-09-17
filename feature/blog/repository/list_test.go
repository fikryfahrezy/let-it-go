package repository_test

import (
	"context"
	"testing"

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
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test pagination
	result, err := testRepository.List(context.Background(), 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 blogs, got %d", len(result))
	}

	result, err = testRepository.List(context.Background(), 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 blogs, got %d", len(result))
	}

	result, err = testRepository.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 blogs, got %d", len(result))
	}
}
