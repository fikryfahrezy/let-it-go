package repository

import (
	"time"
)

type Blog struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Content     string    `db:"content"`
	AuthorID    int       `db:"author_id"`
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