package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	createdBlog, err := getBlogByTitle(blog.Title)
	assert.NoError(t, err)

	err = testRepository.Delete(context.Background(), createdBlog.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = testRepository.GetByID(context.Background(), createdBlog.ID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)
}

func TestDeleteNotFound(t *testing.T) {
	setupTest(t)

	randomID := uuid.New()
	err := testRepository.Delete(context.Background(), randomID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)
}
