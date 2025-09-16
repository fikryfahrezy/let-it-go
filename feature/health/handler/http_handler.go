package handler

import (
	"log/slog"
	"net/http"

	"github.com/fikryfahrezy/let-it-go/pkg/database"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) HealthCheck(c echo.Context) error {
	status := "ok"
	message := "Service is healthy"
	
	// Check database connection
	if err := h.db.Ping(); err != nil {
		slog.Error("Database health check failed",
			slog.String("error", err.Error()),
		)
		status = "unhealthy"
		message = "Database connection failed"
		
		return c.JSON(http.StatusServiceUnavailable, map[string]any{
			"status":  status,
			"message": message,
			"checks": map[string]any{
				"database": map[string]any{
					"status": "unhealthy",
					"error":  err.Error(),
				},
			},
		})
	}
	
	return c.JSON(http.StatusOK, map[string]any{
		"status":  status,
		"message": message,
		"checks": map[string]any{
			"database": map[string]any{
				"status": "healthy",
			},
		},
	})
}

// SetupRoutes configures health check routes
func (h *HealthHandler) SetupRoutes(api *echo.Group) {
	// Health check endpoint (no versioning needed)
	api.GET("/health", h.HealthCheck)
}