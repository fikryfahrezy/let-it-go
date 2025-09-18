package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	// Verify blog was created by getting all blogs and checking the title
	blogs, err := testRepository.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, blogs, 1)

	result := blogs[0]
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, blog.Title, result.Title)
	assert.Equal(t, blog.Content, result.Content)
	assert.Equal(t, blog.AuthorID, result.AuthorID)
	assert.Equal(t, blog.Status, result.Status)
	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())
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
	assert.NoError(t, err)

	result, err := getBlogByTitle(blog.Title)
	assert.NoError(t, err)
	assert.Equal(t, repository.StatusPublished, result.Status)
	assert.NotNil(t, result.PublishedAt)
}
