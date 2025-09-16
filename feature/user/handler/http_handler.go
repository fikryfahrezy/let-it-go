package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req service.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return BadRequestResponse(c, "Invalid request format", err)
	}

	if err := h.validateCreateUserRequest(req); err != nil {
		slog.Warn("Invalid create user request",
			slog.String("error", err.Error()),
		)
		return BadRequestResponse(c, "Validation failed", err)
	}

	user, err := h.userService.CreateUser(c.Request().Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return BadRequestResponse(c, "User creation failed", err)
		}
		slog.Error("Failed to create user",
			slog.String("error", err.Error()),
		)
		return InternalServerErrorResponse(c, "Failed to create user", err)
	}

	return CreatedResponse(c, "User created successfully", user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return BadRequestResponse(c, "Invalid user ID", err)
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return NotFoundResponse(c, "User not found", err)
		}
		slog.Error("Failed to get user",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return InternalServerErrorResponse(c, "Failed to get user", err)
	}

	return SuccessResponse(c, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return BadRequestResponse(c, "Invalid user ID", err)
	}

	var req service.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return BadRequestResponse(c, "Invalid request format", err)
	}

	if err := h.validateUpdateUserRequest(req); err != nil {
		slog.Warn("Invalid update user request",
			slog.String("error", err.Error()),
		)
		return BadRequestResponse(c, "Validation failed", err)
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return NotFoundResponse(c, "User not found", err)
		}
		if strings.Contains(err.Error(), "already exists") {
			return BadRequestResponse(c, "User update failed", err)
		}
		slog.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return InternalServerErrorResponse(c, "Failed to update user", err)
	}

	return SuccessResponse(c, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return BadRequestResponse(c, "Invalid user ID", err)
	}

	err = h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return NotFoundResponse(c, "User not found", err)
		}
		slog.Error("Failed to delete user",
			slog.String("error", err.Error()),
			slog.Int("user_id", id),
		)
		return InternalServerErrorResponse(c, "Failed to delete user", err)
	}

	return SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	pageParam := c.QueryParam("page")
	pageSizeParam := c.QueryParam("page_size")

	page := 1
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	users, pagination, err := h.userService.ListUsers(c.Request().Context(), page, pageSize)
	if err != nil {
		slog.Error("Failed to list users",
			slog.String("error", err.Error()),
		)
		return InternalServerErrorResponse(c, "Failed to list users", err)
	}

	return ListSuccessResponse(c, "Users retrieved successfully", users, pagination)
}

func (h *UserHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "ok",
		"message": "Service is healthy",
	})
}

func (h *UserHandler) validateCreateUserRequest(req service.CreateUserRequest) error {
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "name must be between 2 and 100 characters")
	}
	if req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}
	if !strings.Contains(req.Email, "@") {
		return echo.NewHTTPError(http.StatusBadRequest, "email must be valid")
	}
	if req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password is required")
	}
	if len(req.Password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "password must be at least 6 characters")
	}
	return nil
}

func (h *UserHandler) validateUpdateUserRequest(req service.UpdateUserRequest) error {
	if req.Name != "" && (len(req.Name) < 2 || len(req.Name) > 100) {
		return echo.NewHTTPError(http.StatusBadRequest, "name must be between 2 and 100 characters")
	}
	if req.Email != "" && !strings.Contains(req.Email, "@") {
		return echo.NewHTTPError(http.StatusBadRequest, "email must be valid")
	}
	return nil
}

// SetupRoutes configures all versioned API routes for users
func (h *UserHandler) SetupRoutes(api *echo.Group) {
	// v1 routes
	v1 := api.Group("/v1")
	h.setupV1Routes(v1)
	
	// v2 routes with enhanced features
	v2 := api.Group("/v2")
	h.setupV2Routes(v2)
}

// setupV1Routes configures v1 API routes for users
func (h *UserHandler) setupV1Routes(v1 *echo.Group) {
	users := v1.Group("/users")
	users.POST("", h.CreateUser)
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}

