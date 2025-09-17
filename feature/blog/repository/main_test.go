package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
)

var (
	db             *sql.DB
	testRepository repository.BlogRepository
)

// Helper function to get a blog by title since Create generates new IDs
func getBlogByTitle(title string) (repository.Blog, error) {
	if testRepository == nil {
		return repository.Blog{}, repository.ErrBlogNotFound
	}

	blogs, err := testRepository.List(context.Background(), 100, 0)
	if err != nil {
		return repository.Blog{}, err
	}

	for _, blog := range blogs {
		if blog.Title == title {
			return blog, nil
		}
	}

	return repository.Blog{}, repository.ErrBlogNotFound
}

func TestMain(m *testing.M) {
	fmt.Println("TestMain starting...")
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		fmt.Println("Skipping integration tests")
		os.Exit(0)
	}

	// Create dockertest pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("Failed to create dockertest pool:", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Start MySQL container
	resource, err := pool.Run("mysql", "8.0", []string{
		"MYSQL_ROOT_PASSWORD=testpass",
		"MYSQL_DATABASE=testdb",
		"MYSQL_USER=testuser",
		"MYSQL_PASSWORD=testpass",
	})
	if err != nil {
		log.Fatal("Failed to start MySQL container:", err)
	}

	resource.Expire(60) // 1 minute

	dsn := fmt.Sprintf("testuser:testpass@(localhost:%s)/testdb?parseTime=true", resource.GetPort("3306/tcp"))

	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// nolint:errcheck
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	runMigrations(dsn)
	testRepository = repository.NewBlogRepository(logger.NewDiscardLogger(), &database.DB{DB: db})

	m.Run()
}

func setupTest(t *testing.T) uuid.UUID {
	if db == nil {
		t.Skip("Test database not initialized - set SKIP_INTEGRATION_TESTS=true to skip integration tests")
	}

	// Clean up before each test
	_, err := db.Exec("DELETE FROM blogs")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test user (author)
	authorID := uuid.New()
	_, err = db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		authorID, "Test Author", "author@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	return authorID
}

func runMigrations(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// nolint:errcheck
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
