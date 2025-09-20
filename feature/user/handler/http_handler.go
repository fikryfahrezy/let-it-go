package handler

import (
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
	log         *slog.Logger
}

func NewUserHandler(log *slog.Logger, userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

// translateServiceError converts service errors to appropriate HTTP responses
func (h *UserHandler) translateServiceError(c echo.Context, err error, defaultMessage string) error {
	if errors.Is(err, service.ErrUserAlreadyExists) {
		return http_server.BadRequestResponse(c, "Email address is already taken", err)
	}
	if errors.Is(err, service.ErrFailedToHashPassword) {
		return http_server.InternalServerErrorResponse(c, "Password processing failed", err)
	}
	if errors.Is(err, repository.ErrUserNotFound) {
		return http_server.NotFoundResponse(c, "User not found", err)
	}

	// Log unexpected errors
	h.log.Error("Service error",
		slog.String("error", err.Error()),
		slog.String("operation", defaultMessage),
	)
	return http_server.InternalServerErrorResponse(c, defaultMessage, err)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body service.CreateUserRequest true "User creation request"
// @Success 201 {object} http_server.APIResponse{result=service.GetUserResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req service.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return http_server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := c.Validate(&req); err != nil {
		return http_server.HandleValidationError(c, err)
	}

	user, err := h.userService.CreateUser(c.Request().Context(), req)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to create user")
	}

	return http_server.CreatedResponse(c, "User created successfully", user)
}

// GetUser retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieve a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} http_server.APIResponse{result=service.GetUserResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid user UUID format", err)
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to get user")
	}

	return http_server.SuccessResponse(c, "User retrieved successfully", user)
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body service.UpdateUserRequest true "User update request"
// @Success 200 {object} http_server.APIResponse{result=service.GetUserResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid user UUID format", err)
	}

	var req service.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return http_server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := c.Validate(&req); err != nil {
		return http_server.HandleValidationError(c, err)
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to update user")
	}

	return http_server.SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser deletes a user by ID
// @Summary Delete a user
// @Description Delete a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} http_server.APIResponse
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid user ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid user UUID format", err)
	}

	err = h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to delete user")
	}

	return http_server.SuccessResponse(c, "User deleted successfully", nil)
}

// ListUsers retrieves a list of users with pagination
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} http_server.ListAPIResponse{result=[]service.GetUserResponse}
// @Failure 500 {object} http_server.APIResponse
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

	paginationReq := http_server.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
	users, totalCount, err := h.userService.ListUsers(c.Request().Context(), service.ListUsersRequest{
		PaginationRequest: paginationReq,
	})
	if err != nil {
		return h.translateServiceError(c, err, "Failed to list users")
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	pagination := http_server.CreatePaginationResponse(totalCount, totalPages, page, pageSize)

	return http_server.ListSuccessResponse(c, "Users retrieved successfully", users, pagination)
}

func (h *UserHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "ok",
		"message": "Service is healthy",
	})
}

// SetupRoutes configures all versioned API routes for users
func (h *UserHandler) SetupRoutes(server *http_server.Server) {
	// v1 routes
	h.setupV1Routes(server)

	// v2 routes with enhanced features
	h.setupV2Routes(server)
}

// setupV1Routes configures v1 API routes for users
func (h *UserHandler) setupV1Routes(server *http_server.Server) {
	users := server.Echo().Group("/v1/users")
	users.POST("", h.CreateUser)
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}
