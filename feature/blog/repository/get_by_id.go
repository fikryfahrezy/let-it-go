package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func (r *blogRepository) GetByID(ctx context.Context, id int) (Blog, error) {
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
			return Blog{}, fmt.Errorf("blog not found")
		}
		slog.Error("Failed to get blog by ID",
			slog.String("error", err.Error()),
			slog.Int("blog_id", id),
		)
		return Blog{}, fmt.Errorf("failed to get blog by ID: %w", err)
	}

	return blog, nil
}