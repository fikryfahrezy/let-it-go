package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestGetByID(t *testing.T) {
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

	result, err := testRepository.GetByID(context.Background(), createdBlog.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != createdBlog.ID {
		t.Errorf("Expected ID %s, got %s", createdBlog.ID, result.ID)
	}
	if result.Title != blog.Title {
		t.Errorf("Expected title %s, got %s", blog.Title, result.Title)
	}
	if result.Content != blog.Content {
		t.Errorf("Expected content %s, got %s", blog.Content, result.Content)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	setupTest(t)

	randomID := uuid.New()
	_, err := testRepository.GetByID(context.Background(), randomID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != repository.ErrBlogNotFound {
		t.Errorf("Expected ErrBlogNotFound, got %v", err)
	}
}
