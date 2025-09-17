package repository

import (
	"context"
	"fmt"
	"log/slog"
)

func (r *blogRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]Blog, error) {
	query := `
		SELECT id, title, content, author_id, status, published_at, created_at, updated_at
		FROM blogs
		WHERE status = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		r.log.Error("Failed to get blogs by status",
			slog.String("error", err.Error()),
			slog.String("status", status),
		)
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetBlogsByStatus, err)
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
			r.log.Error("Failed to scan blog row",
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("%w: %w", ErrFailedToScanBlogRow, err)
		}
		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("Error iterating blog rows",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%w: %w", ErrFailedToIterateRows, err)
	}

	return blogs, nil
}
