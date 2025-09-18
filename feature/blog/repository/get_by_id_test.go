package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, createdBlog.ID, result.ID)
	assert.Equal(t, blog.Title, result.Title)
	assert.Equal(t, blog.Content, result.Content)
}

func TestGetByIDNotFound(t *testing.T) {
	setupTest(t)

	randomID := uuid.New()
	_, err := testRepository.GetByID(context.Background(), randomID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)
}
