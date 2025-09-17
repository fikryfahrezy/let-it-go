package repository_test

import (
	"context"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestCount(t *testing.T) {
	authorID := setupTest(t)

	// Initially empty
	count, err := testRepository.Count(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

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
		if err != nil {
			t.Fatal(err)
		}
	}

	count, err = testRepository.Count(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

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
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test count by status
	count, err := testRepository.CountByStatus(context.Background(), repository.StatusDraft)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("Expected draft count 2, got %d", count)
	}

	count, err = testRepository.CountByStatus(context.Background(), repository.StatusPublished)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected published count 1, got %d", count)
	}

	count, err = testRepository.CountByStatus(context.Background(), repository.StatusArchived)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("Expected archived count 0, got %d", count)
	}
}
