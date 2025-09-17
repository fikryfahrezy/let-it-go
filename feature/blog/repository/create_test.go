package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestCreate(t *testing.T) {
	authorID := setupTest(t)

	blog := repository.Blog{
		Title:    "Test Blog",
		Content:  "This is a test blog content",
		AuthorID: authorID,
		Status:   repository.StatusDraft,
	}

	err := testRepository.Create(context.Background(), blog)
	if err != nil {
		t.Fatal(err)
	}

	// Verify blog was created by getting all blogs and checking the title
	blogs, err := testRepository.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(blogs) != 1 {
		t.Errorf("Expected 1 blog, got %d", len(blogs))
	}

	result := blogs[0]
	if result.ID == uuid.Nil {
		t.Error("Expected non-nil ID")
	}
	if result.Title != blog.Title {
		t.Errorf("Expected title %s, got %s", blog.Title, result.Title)
	}
	if result.Content != blog.Content {
		t.Errorf("Expected content %s, got %s", blog.Content, result.Content)
	}
	if result.AuthorID != blog.AuthorID {
		t.Errorf("Expected authorID %s, got %s", blog.AuthorID, result.AuthorID)
	}
	if result.Status != blog.Status {
		t.Errorf("Expected status %s, got %s", blog.Status, result.Status)
	}
	if result.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
	if result.UpdatedAt.IsZero() {
		t.Error("Expected non-zero UpdatedAt")
	}
}

func TestCreateWithPublishedStatus(t *testing.T) {
	authorID := setupTest(t)

	publishedAt := time.Now().UTC()
	blog := repository.Blog{
		Title:       "Published Blog",
		Content:     "This is a published blog content",
		AuthorID:    authorID,
		Status:      repository.StatusPublished,
		PublishedAt: &publishedAt,
	}

	err := testRepository.Create(context.Background(), blog)
	if err != nil {
		t.Fatal(err)
	}

	result, err := getBlogByTitle(blog.Title)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != repository.StatusPublished {
		t.Errorf("Expected status %s, got %s", repository.StatusPublished, result.Status)
	}
	if result.PublishedAt == nil {
		t.Error("Expected non-nil PublishedAt")
	}
}
