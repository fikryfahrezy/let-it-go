-- Migration: create_blogs_table
-- Created: 2025-01-16T24:47:06Z

-- Create blogs table
CREATE TABLE IF NOT EXISTS blogs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_author_id (author_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_published_at (published_at),
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Insert sample data
INSERT IGNORE INTO blogs (title, content, author_id, status, published_at) VALUES 
('Welcome to Our Blog', 'This is our first blog post. Welcome to our amazing blog where we share insights about technology, development, and more!', 1, 'published', NOW()),
('Draft Post', 'This is a draft post that is not yet published. It contains some preliminary ideas that need further development.', 1, 'draft', NULL),
('Getting Started with Go', 'Go is a powerful programming language that makes it easy to build simple, reliable, and efficient software. In this post, we will explore the basics of Go programming.', 2, 'published', NOW() - INTERVAL 1 DAY),
('Advanced Go Patterns', 'In this post, we dive deep into advanced Go programming patterns including interfaces, goroutines, and channels.', 2, 'draft', NULL);