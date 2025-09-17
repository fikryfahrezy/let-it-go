package repository_test

import (
	"context"
	"testing"

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
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test get by draft status
	result, err := testRepository.GetByStatus(context.Background(), repository.StatusDraft, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 draft blogs, got %d", len(result))
	}

	// Test get by published status
	result, err = testRepository.GetByStatus(context.Background(), repository.StatusPublished, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 published blog, got %d", len(result))
	}

	// Test get by archived status
	result, err = testRepository.GetByStatus(context.Background(), repository.StatusArchived, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 archived blog, got %d", len(result))
	}

	// Test pagination
	result, err = testRepository.GetByStatus(context.Background(), repository.StatusDraft, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 draft blog with pagination, got %d", len(result))
	}
}
