package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
)

func TestDelete(t *testing.T) {
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

	err = testRepository.Delete(context.Background(), createdBlog.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify deletion
	_, err = testRepository.GetByID(context.Background(), createdBlog.ID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != repository.ErrBlogNotFound {
		t.Errorf("Expected ErrBlogNotFound, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	setupTest(t)

	randomID := uuid.New()
	err := testRepository.Delete(context.Background(), randomID)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != repository.ErrBlogNotFound {
		t.Errorf("Expected ErrBlogNotFound, got %v", err)
	}
}
