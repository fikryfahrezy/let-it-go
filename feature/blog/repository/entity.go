package repository

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID          uuid.UUID `db:"id"` // UUIDv7
	Title       string    `db:"title"`
	Content     string    `db:"content"`
	AuthorID    uuid.UUID `db:"author_id"` // UUIDv7
	Status      string    `db:"status"`
	PublishedAt *time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)