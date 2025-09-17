package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type BlogRepositoryTestSuite struct {
	suite.Suite
	container  *mysqlcontainer.MySQLContainer
	db         *database.DB
	repository BlogRepository
	ctx        context.Context
	authorID   uuid.UUID
}

// Helper function to get a blog by title since Create generates new IDs
func (suite *BlogRepositoryTestSuite) getBlogByTitle(title string) (Blog, error) {
	blogs, err := suite.repository.List(suite.ctx, 100, 0)
	if err != nil {
		return Blog{}, err
	}

	for _, blog := range blogs {
		if blog.Title == title {
			return blog, nil
		}
	}

	return Blog{}, ErrBlogNotFound
}

func (suite *BlogRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Start MySQL container
	container, err := mysqlcontainer.Run(suite.ctx,
		"mysql:8.0",
		mysqlcontainer.WithDatabase("testdb"),
		mysqlcontainer.WithUsername("testuser"),
		mysqlcontainer.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server - GPL").
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Fatal("Failed to start MySQL container:", err)
	}
	suite.container = container

	// Get connection string
	dsn, err := container.ConnectionString(suite.ctx, "parseTime=true")
	if err != nil {
		log.Fatal("Failed to get connection string:", err)
	}

	// Connect to database
	suite.db, err = database.NewDB(database.Config{DSN: dsn})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	suite.runMigrations(dsn)

	// Initialize repository
	suite.repository = NewBlogRepository(suite.db)
}

func (suite *BlogRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(suite.ctx)
	}
}

func (suite *BlogRepositoryTestSuite) SetupTest() {
	// Clean up before each test
	_, err := suite.db.Exec("DELETE FROM blogs")
	require.NoError(suite.T(), err)
	_, err = suite.db.Exec("DELETE FROM users")
	require.NoError(suite.T(), err)

	// Create a test user (author)
	suite.authorID = uuid.New()
	_, err = suite.db.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		suite.authorID, "Test Author", "author@example.com", "password")
	require.NoError(suite.T(), err)
}

func (suite *BlogRepositoryTestSuite) runMigrations(dsn string) {
	db, err := sql.Open("mysql", dsn)
	require.NoError(suite.T(), err)
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	require.NoError(suite.T(), err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"mysql",
		driver,
	)
	require.NoError(suite.T(), err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(suite.T(), err)
	}
}
