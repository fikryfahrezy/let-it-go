package repository_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
)

var (
	db             *database.DB
	testRepository repository.UserRepository
)

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
		db, err = database.NewDB(database.Config{DSN: dsn})
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
	testRepository = repository.NewUserRepository(logger.NewDiscardLogger(), db)

	m.Run()
}

func setupTest(t *testing.T) {
	if db == nil {
		t.Skip("Test database not initialized - set SKIP_INTEGRATION_TESTS=true to skip integration tests")
	}

	// Clean up before each test
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatal(err)
	}
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
