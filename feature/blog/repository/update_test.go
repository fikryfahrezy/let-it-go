package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestUpdate(t *testing.T) {
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

	createdBlog, err := getBlogByTitle(blog.Title)
	if err != nil {
		t.Fatal(err)
	}

	// Update the blog
	publishedAt := time.Now().UTC()
	updatedBlog := repository.Blog{
		ID:          createdBlog.ID,
		Title:       "Updated Test Blog",
		Content:     "This is updated test blog content",
		AuthorID:    authorID,
		Status:      repository.StatusPublished,
		PublishedAt: &publishedAt,
	}

	err = testRepository.Update(context.Background(), updatedBlog)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the update
	result, err := testRepository.GetByID(context.Background(), createdBlog.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Title != updatedBlog.Title {
		t.Errorf("Expected title %s, got %s", updatedBlog.Title, result.Title)
	}
	if result.Content != updatedBlog.Content {
		t.Errorf("Expected content %s, got %s", updatedBlog.Content, result.Content)
	}
	if result.Status != updatedBlog.Status {
		t.Errorf("Expected status %s, got %s", updatedBlog.Status, result.Status)
	}
	if result.PublishedAt == nil {
		t.Error("Expected non-nil PublishedAt")
	}
}

func TestUpdateNotFound(t *testing.T) {
	setupTest(t)

	randomID := uuid.New()
	blog := repository.Blog{
		ID:      randomID,
		Title:   "Non-existent Blog",
		Content: "Content",
		Status:  repository.StatusDraft,
	}

	err := testRepository.Update(context.Background(), blog)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != repository.ErrBlogNotFound {
		t.Errorf("Expected ErrBlogNotFound, got %v", err)
	}
}
