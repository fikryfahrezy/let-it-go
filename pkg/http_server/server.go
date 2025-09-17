package http_server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Config holds the server configuration
type Config struct {
	Host string
	Port int
}

// RouteHandler defines interface for features to register their routes
type RouteHandler interface {
	SetupRoutes(api *echo.Group)
}

// Server manages the application server lifecycle
type Server struct {
	config Config
	echo   *echo.Echo
}

// New creates a new server instance
func New(config Config) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}

// Initialize sets up the server with dependencies
func (s *Server) Initialize(handlers []RouteHandler) error {
	slog.Info("Initializing server")

	// Initialize Echo server
	e := echo.New()

	// Set custom validator
	e.Validator = NewCustomValidator()

	// Configure middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Request logging with slog
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("HTTP Request",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
			)
			return nil
		},
	}))

	// Setup API routes
	s.setupAPIRoutes(e, handlers)

	s.echo = e
	return nil
}

// setupAPIRoutes configures all API routes for all versions
func (s *Server) setupAPIRoutes(e *echo.Echo, handlers []RouteHandler) {
	// Swagger documentation - accessible at /swagger/index.html
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Redirect /swagger to /swagger/index.html for convenience
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(302, "/swagger/index.html")
	})

	api := e.Group("/api")

	// Let each feature register its own routes
	for _, handler := range handlers {
		handler.SetupRoutes(api)
	}
}

// Start starts the server
func (s *Server) Start() error {
	if s.echo == nil {
		return fmt.Errorf("server not initialized, call Initialize() first")
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	slog.Info("Starting server",
		slog.String("address", addr),
	)
	return s.echo.Start(addr)
}

// Stop stops the server gracefully
func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stopping server...")
	return s.echo.Shutdown(ctx)
}
