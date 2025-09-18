package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	createdBlog, err := getBlogByTitle(blog.Title)
	assert.NoError(t, err)

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
	assert.NoError(t, err)

	// Verify the update
	result, err := testRepository.GetByID(context.Background(), createdBlog.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedBlog.Title, result.Title)
	assert.Equal(t, updatedBlog.Content, result.Content)
	assert.Equal(t, updatedBlog.Status, result.Status)
	assert.NotNil(t, result.PublishedAt)
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
	assert.Error(t, err)
	assert.Equal(t, repository.ErrBlogNotFound, err)
}
