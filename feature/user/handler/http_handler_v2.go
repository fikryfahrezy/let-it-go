package handler

import (
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/labstack/echo/v4"
)

// setupV2Routes configures v2 API routes for users with enhanced features
func (h *UserHandler) setupV2Routes(v2 *echo.Group) {
	users := v2.Group("/users")
	users.POST("", h.CreateUser)
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
	
	// v2 specific endpoints
	users.GET("/:id/profile", h.GetUserProfile)
	users.POST("/batch", h.BatchUserOperations)
}

// GetUserProfile handles enhanced user profile endpoint for v2
func (h *UserHandler) GetUserProfile(c echo.Context) error {
	data := map[string]any{
		"user_id": c.Param("id"),
		"version": "v2",
	}
	return http_server.SuccessResponse(c, "Enhanced user profile endpoint (v2)", data)
}

// BatchUserOperations handles batch user operations for v2
func (h *UserHandler) BatchUserOperations(c echo.Context) error {
	data := map[string]any{
		"version": "v2",
		"status": "pending",
	}
	return http_server.SuccessResponse(c, "Batch user operations (v2)", data)
}