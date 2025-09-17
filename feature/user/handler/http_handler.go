package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
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

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body service.CreateUserRequest true "User creation request"
// @Success 201 {object} APIResponse{data=service.GetUserResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req service.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := h.validateCreateUserRequest(req); err != nil {
		slog.Warn("Invalid create user request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Validation failed", err)
	}

	user, err := h.userService.CreateUser(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return server.BadRequestResponse(c, "Email address is already taken", err)
		}
		if errors.Is(err, service.ErrFailedToHashPassword) {
			return server.InternalServerErrorResponse(c, "Password processing failed", err)
		}
		slog.Error("Failed to create user",
			slog.String("error", err.Error()),
		)
		return server.InternalServerErrorResponse(c, "Failed to create user", err)
	}

	return server.CreatedResponse(c, "User created successfully", user)
}

// GetUser retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieve a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} APIResponse{data=service.GetUserResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid user UUID format", err)
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return server.NotFoundResponse(c, "User not found", err)
		}
		slog.Error("Failed to get user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to get user", err)
	}

	return server.SuccessResponse(c, "User retrieved successfully", user)
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body service.UpdateUserRequest true "User update request"
// @Success 200 {object} APIResponse{data=service.GetUserResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid user UUID format", err)
	}

	var req service.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := h.validateUpdateUserRequest(req); err != nil {
		slog.Warn("Invalid update user request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Validation failed", err)
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return server.NotFoundResponse(c, "User not found", err)
		}
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return server.BadRequestResponse(c, "Email address is already taken", err)
		}
		slog.Error("Failed to update user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to update user", err)
	}

	return server.SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser deletes a user by ID
// @Summary Delete a user
// @Description Delete a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid user UUID format", err)
	}

	err = h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return server.NotFoundResponse(c, "User not found", err)
		}
		slog.Error("Failed to delete user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to delete user", err)
	}

	return server.SuccessResponse(c, "User deleted successfully", nil)
}

// ListUsers retrieves a list of users with pagination
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} ListAPIResponse{data=[]service.GetUserResponse}
// @Failure 500 {object} APIResponse
// @Router /v1/users [get]
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
		return server.InternalServerErrorResponse(c, "Failed to list users", err)
	}

	return server.ListSuccessResponse(c, "Users retrieved successfully", users, pagination)
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
