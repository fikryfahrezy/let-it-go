package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (r *blogRepository) GetByID(ctx context.Context, id uuid.UUID) (Blog, error) {
	query := `
		SELECT id, title, content, author_id, status, published_at, created_at, updated_at
		FROM blogs
		WHERE id = ?
	`

	var blog Blog
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.AuthorID,
		&blog.Status,
		&blog.PublishedAt,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Blog{}, ErrBlogNotFound
		}
		r.log.Error("Failed to get blog by ID",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return Blog{}, fmt.Errorf("%w: %w", ErrFailedToGetBlog, err)
	}

	return blog, nil
}
