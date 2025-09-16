package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *blogRepository) List(ctx context.Context, limit, offset int) ([]Blog, error) {
	query := `
		SELECT id, title, content, author_id, status, published_at, created_at, updated_at
		FROM blogs
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Failed to list blogs",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to list blogs: %w", err)
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		blog := Blog{}
		err := rows.Scan(
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
			slog.Error("Failed to scan blog row",
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("failed to scan blog row: %w", err)
		}
		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating blog rows",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("error iterating blog rows: %w", err)
	}

	return blogs, nil
}
