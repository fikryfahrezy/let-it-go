package handler

import (
	"log/slog"
	"net/http"

	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	db        *database.DB
	version   string
	commit    string
	buildTime string
}

func NewHealthHandler(db *database.DB, version, commit, buildTime string) *HealthHandler {
	return &HealthHandler{
		db:        db,
		version:   version,
		commit:    commit,
		buildTime: buildTime,
	}
}

func (h *HealthHandler) HealthCheck(c echo.Context) error {
	status := "ok"
	message := "Service is healthy"
	httpStatus := http.StatusOK

	dbCheck := map[string]any{"status": "healthy"}

	// Check database connection
	if err := h.db.Ping(); err != nil {
		slog.Error("Database health check failed",
			slog.String("error", err.Error()),
		)
		status = "unhealthy"
		message = "Database connection failed"
		httpStatus = http.StatusServiceUnavailable
		dbCheck = map[string]any{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	}

	response := map[string]any{
		"status":    status,
		"message":   message,
		"version":   h.version,
		"commit":    h.commit,
		"buildTime": h.buildTime,
		"checks": map[string]any{
			"database": dbCheck,
		},
	}

	return c.JSON(httpStatus, response)
}

// SetupRoutes configures health check routes
func (h *HealthHandler) SetupRoutes(api *echo.Group) {
	// Health check endpoint (no versioning needed)
	api.GET("/health", h.HealthCheck)
}
