package handler

import (
	"errors"
	"log/slog"
	"math"
	"strconv"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BlogHandler struct {
	blogService service.BlogService
	log         *slog.Logger
}

func NewBlogHandler(log *slog.Logger, blogService service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
		log:         log,
	}
}

// translateServiceError converts service errors to appropriate HTTP responses
func (h *BlogHandler) translateServiceError(c echo.Context, err error, defaultMessage string) error {
	if errors.Is(err, repository.ErrBlogNotFound) {
		return http_server.NotFoundResponse(c, "Blog not found", err)
	}
	if errors.Is(err, service.ErrInvalidBlogStatus) {
		return http_server.BadRequestResponse(c, "Invalid blog status", err)
	}
	if errors.Is(err, service.ErrBlogAlreadyPublished) {
		return http_server.BadRequestResponse(c, "Blog is already published", err)
	}
	if errors.Is(err, service.ErrBlogAlreadyArchived) {
		return http_server.BadRequestResponse(c, "Blog is already archived", err)
	}
	if errors.Is(err, service.ErrFailedToPublishBlog) {
		return http_server.InternalServerErrorResponse(c, "Failed to publish blog", err)
	}
	if errors.Is(err, service.ErrFailedToArchiveBlog) {
		return http_server.InternalServerErrorResponse(c, "Failed to archive blog", err)
	}

	// Log unexpected errors
	h.log.Error("Service error",
		slog.String("error", err.Error()),
		slog.String("operation", defaultMessage),
	)
	return http_server.InternalServerErrorResponse(c, defaultMessage, err)
}

// CreateBlog creates a new blog
// @Summary Create a new blog
// @Description Create a new blog with the provided information
// @Tags blogs
// @Accept json
// @Produce json
// @Param blog body service.CreateBlogRequest true "Blog creation request"
// @Success 201 {object} http_server.APIResponse{result=service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs [post]
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	var req service.CreateBlogRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return http_server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := c.Validate(&req); err != nil {
		return http_server.HandleValidationError(c, err)
	}

	blog, err := h.blogService.CreateBlog(c.Request().Context(), req)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to create blog")
	}

	return http_server.CreatedResponse(c, "Blog created successfully", blog)
}

// GetBlog retrieves a blog by ID
// @Summary Get a blog by ID
// @Description Retrieve a blog by its unique identifier
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} http_server.APIResponse{result=service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/{id} [get]
func (h *BlogHandler) GetBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	blog, err := h.blogService.GetBlogByID(c.Request().Context(), id)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to get blog")
	}

	return http_server.SuccessResponse(c, "Blog retrieved successfully", blog)
}

// UpdateBlog updates an existing blog
// @Summary Update a blog
// @Description Update an existing blog with the provided information
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param blog body service.UpdateBlogRequest true "Blog update request"
// @Success 200 {object} http_server.APIResponse{result=service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/{id} [put]
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	var req service.UpdateBlogRequest
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return http_server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := c.Validate(&req); err != nil {
		return http_server.HandleValidationError(c, err)
	}

	blog, err := h.blogService.UpdateBlog(c.Request().Context(), id, req)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to update blog")
	}

	return http_server.SuccessResponse(c, "Blog updated successfully", blog)
}

// DeleteBlog deletes a blog by ID
// @Summary Delete a blog
// @Description Delete a blog by its unique identifier
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} http_server.APIResponse
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/{id} [delete]
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	err = h.blogService.DeleteBlog(c.Request().Context(), id)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to delete blog")
	}

	return http_server.SuccessResponse(c, "Blog deleted successfully", nil)
}

// ListBlogs retrieves a list of blogs with pagination
// @Summary List blogs
// @Description Retrieve a paginated list of blogs
// @Tags blogs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} http_server.ListAPIResponse{result=[]service.GetBlogResponse}
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs [get]
func (h *BlogHandler) ListBlogs(c echo.Context) error {
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
	blogs, totalCount, err := h.blogService.ListBlogs(c.Request().Context(), service.ListBlogsRequest{
		PaginationRequest: paginationReq,
	})
	if err != nil {
		return h.translateServiceError(c, err, "Failed to list blogs")
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	pagination := http_server.CreatePaginationResponse(totalCount, totalPages, page, pageSize)

	return http_server.ListSuccessResponse(c, "Blogs retrieved successfully", blogs, pagination)
}

// GetBlogsByAuthor retrieves blogs by author ID with pagination
// @Summary Get blogs by author
// @Description Retrieve a paginated list of blogs by author ID
// @Tags blogs
// @Accept json
// @Produce json
// @Param author_id path string true "Author ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} http_server.ListAPIResponse{result=[]service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/author/{author_id} [get]
func (h *BlogHandler) GetBlogsByAuthor(c echo.Context) error {
	authorIDParam := c.Param("author_id")
	authorID, err := uuid.Parse(authorIDParam)
	if err != nil {
		h.log.Warn("Invalid author ID parameter",
			slog.String("author_id", authorIDParam),
		)
		return http_server.BadRequestResponse(c, "Invalid author UUID format", err)
	}

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
	blogs, totalCount, err := h.blogService.GetBlogsByAuthor(c.Request().Context(), authorID, service.GetBlogsByAuthorRequest{
		PaginationRequest: paginationReq,
	})
	if err != nil {
		return h.translateServiceError(c, err, "Failed to get blogs by author")
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	pagination := http_server.CreatePaginationResponse(totalCount, totalPages, page, pageSize)

	return http_server.ListSuccessResponse(c, "Blogs retrieved successfully", blogs, pagination)
}

// GetBlogsByStatus retrieves blogs by status with pagination
// @Summary Get blogs by status
// @Description Retrieve a paginated list of blogs by status
// @Tags blogs
// @Accept json
// @Produce json
// @Param status path string true "Blog status" Enums(draft, published, archived)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} http_server.ListAPIResponse{result=[]service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/status/{status} [get]
func (h *BlogHandler) GetBlogsByStatus(c echo.Context) error {
	status := c.Param("status")
	if status == "" {
		return http_server.BadRequestResponse(c, "Status is required", nil)
	}

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
	blogs, totalCount, err := h.blogService.GetBlogsByStatus(c.Request().Context(), status, service.GetBlogsByStatusRequest{
		PaginationRequest: paginationReq,
	})
	if err != nil {
		return h.translateServiceError(c, err, "Failed to get blogs by status")
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	pagination := http_server.CreatePaginationResponse(totalCount, totalPages, page, pageSize)

	return http_server.ListSuccessResponse(c, "Blogs retrieved successfully", blogs, pagination)
}

// PublishBlog publishes a blog by ID
// @Summary Publish a blog
// @Description Publish a blog by setting its status to published
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} http_server.APIResponse{result=service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/{id}/publish [post]
func (h *BlogHandler) PublishBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	blog, err := h.blogService.PublishBlog(c.Request().Context(), id)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to publish blog")
	}

	return http_server.SuccessResponse(c, "Blog published successfully", blog)
}

// ArchiveBlog archives a blog by ID
// @Summary Archive a blog
// @Description Archive a blog by setting its status to archived
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} http_server.APIResponse{result=service.GetBlogResponse}
// @Failure 400 {object} http_server.APIResponse
// @Failure 404 {object} http_server.APIResponse
// @Failure 500 {object} http_server.APIResponse
// @Router /v1/blogs/{id}/archive [post]
func (h *BlogHandler) ArchiveBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return http_server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	blog, err := h.blogService.ArchiveBlog(c.Request().Context(), id)
	if err != nil {
		return h.translateServiceError(c, err, "Failed to archive blog")
	}

	return http_server.SuccessResponse(c, "Blog archived successfully", blog)
}

// SetupRoutes configures all API routes for blogs
func (h *BlogHandler) SetupRoutes(server *http_server.Server) {
	h.setupV1Routes(server)
}

// setupV1Routes configures v1 API routes for blogs
func (h *BlogHandler) setupV1Routes(server *http_server.Server) {
	blogs := server.Echo().Group("/v1/blogs")
	blogs.POST("", h.CreateBlog)
	blogs.GET("", h.ListBlogs)
	blogs.GET("/:id", h.GetBlog)
	blogs.PUT("/:id", h.UpdateBlog)
	blogs.DELETE("/:id", h.DeleteBlog)
	blogs.GET("/author/:author_id", h.GetBlogsByAuthor)
	blogs.GET("/status/:status", h.GetBlogsByStatus)
	blogs.POST("/:id/publish", h.PublishBlog)
	blogs.POST("/:id/archive", h.ArchiveBlog)
}
