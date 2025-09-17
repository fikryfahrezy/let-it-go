package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BlogHandler struct {
	blogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

// CreateBlog creates a new blog
// @Summary Create a new blog
// @Description Create a new blog with the provided information
// @Tags blogs
// @Accept json
// @Produce json
// @Param blog body service.CreateBlogRequest true "Blog creation request"
// @Success 201 {object} APIResponse{data=service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs [post]
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	var req service.CreateBlogRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := h.validateCreateBlogRequest(req); err != nil {
		slog.Warn("Invalid create blog request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Validation failed", err)
	}

	blog, err := h.blogService.CreateBlog(c.Request().Context(), req)
	if err != nil {
		slog.Error("Failed to create blog",
			slog.String("error", err.Error()),
		)
		return server.InternalServerErrorResponse(c, "Failed to create blog", err)
	}

	return server.CreatedResponse(c, "Blog created successfully", blog)
}

// GetBlog retrieves a blog by ID
// @Summary Get a blog by ID
// @Description Retrieve a blog by its unique identifier
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} APIResponse{data=service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/{id} [get]
func (h *BlogHandler) GetBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	blog, err := h.blogService.GetBlogByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return server.NotFoundResponse(c, "Blog not found", err)
		}
		slog.Error("Failed to get blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to get blog", err)
	}

	return server.SuccessResponse(c, "Blog retrieved successfully", blog)
}

// UpdateBlog updates an existing blog
// @Summary Update a blog
// @Description Update an existing blog with the provided information
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param blog body service.UpdateBlogRequest true "Blog update request"
// @Success 200 {object} APIResponse{data=service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/{id} [put]
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	var req service.UpdateBlogRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Invalid request format", err)
	}

	if err := h.validateUpdateBlogRequest(req); err != nil {
		slog.Warn("Invalid update blog request",
			slog.String("error", err.Error()),
		)
		return server.BadRequestResponse(c, "Validation failed", err)
	}

	blog, err := h.blogService.UpdateBlog(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return server.NotFoundResponse(c, "Blog not found", err)
		}
		slog.Error("Failed to update blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to update blog", err)
	}

	return server.SuccessResponse(c, "Blog updated successfully", blog)
}

// DeleteBlog deletes a blog by ID
// @Summary Delete a blog
// @Description Delete a blog by its unique identifier
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/{id} [delete]
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	err = h.blogService.DeleteBlog(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return server.NotFoundResponse(c, "Blog not found", err)
		}
		slog.Error("Failed to delete blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to delete blog", err)
	}

	return server.SuccessResponse(c, "Blog deleted successfully", nil)
}

// ListBlogs retrieves a list of blogs with pagination
// @Summary List blogs
// @Description Retrieve a paginated list of blogs
// @Tags blogs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} ListAPIResponse{data=[]service.GetBlogResponse}
// @Failure 500 {object} APIResponse
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

	blogs, pagination, err := h.blogService.ListBlogs(c.Request().Context(), page, pageSize)
	if err != nil {
		slog.Error("Failed to list blogs",
			slog.String("error", err.Error()),
		)
		return server.InternalServerErrorResponse(c, "Failed to list blogs", err)
	}

	return server.ListSuccessResponse(c, "Blogs retrieved successfully", blogs, pagination)
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
// @Success 200 {object} ListAPIResponse{data=[]service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/author/{author_id} [get]
func (h *BlogHandler) GetBlogsByAuthor(c echo.Context) error {
	authorIDParam := c.Param("author_id")
	authorID, err := uuid.Parse(authorIDParam)
	if err != nil {
		slog.Warn("Invalid author ID parameter",
			slog.String("author_id", authorIDParam),
		)
		return server.BadRequestResponse(c, "Invalid author UUID format", err)
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

	blogs, pagination, err := h.blogService.GetBlogsByAuthor(c.Request().Context(), authorID, page, pageSize)
	if err != nil {
		slog.Error("Failed to get blogs by author",
			slog.String("error", err.Error()),
			slog.String("author_id", authorID.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to get blogs by author", err)
	}

	return server.ListSuccessResponse(c, "Blogs retrieved successfully", blogs, pagination)
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
// @Success 200 {object} ListAPIResponse{data=[]service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/status/{status} [get]
func (h *BlogHandler) GetBlogsByStatus(c echo.Context) error {
	status := c.Param("status")
	if status == "" {
		return server.BadRequestResponse(c, "Status is required", nil)
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

	blogs, pagination, err := h.blogService.GetBlogsByStatus(c.Request().Context(), status, page, pageSize)
	if err != nil {
		slog.Error("Failed to get blogs by status",
			slog.String("error", err.Error()),
			slog.String("status", status),
		)
		return server.InternalServerErrorResponse(c, "Failed to get blogs by status", err)
	}

	return server.ListSuccessResponse(c, "Blogs retrieved successfully", blogs, pagination)
}

// PublishBlog publishes a blog by ID
// @Summary Publish a blog
// @Description Publish a blog by setting its status to published
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} APIResponse{data=service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/{id}/publish [post]
func (h *BlogHandler) PublishBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	blog, err := h.blogService.PublishBlog(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return server.NotFoundResponse(c, "Blog not found", err)
		}
		slog.Error("Failed to publish blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to publish blog", err)
	}

	return server.SuccessResponse(c, "Blog published successfully", blog)
}

// ArchiveBlog archives a blog by ID
// @Summary Archive a blog
// @Description Archive a blog by setting its status to archived
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} APIResponse{data=service.GetBlogResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /v1/blogs/{id}/archive [post]
func (h *BlogHandler) ArchiveBlog(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("Invalid blog ID parameter",
			slog.String("id", idParam),
		)
		return server.BadRequestResponse(c, "Invalid blog UUID format", err)
	}

	blog, err := h.blogService.ArchiveBlog(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBlogNotFound) {
			return server.NotFoundResponse(c, "Blog not found", err)
		}
		slog.Error("Failed to archive blog",
			slog.String("error", err.Error()),
			slog.String("blog_id", id.String()),
		)
		return server.InternalServerErrorResponse(c, "Failed to archive blog", err)
	}

	return server.SuccessResponse(c, "Blog archived successfully", blog)
}

func (h *BlogHandler) validateCreateBlogRequest(req service.CreateBlogRequest) error {
	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title is required")
	}
	if len(req.Title) < 2 || len(req.Title) > 200 {
		return echo.NewHTTPError(http.StatusBadRequest, "title must be between 2 and 200 characters")
	}
	if req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "content is required")
	}
	if len(req.Content) < 10 {
		return echo.NewHTTPError(http.StatusBadRequest, "content must be at least 10 characters")
	}
	if req.AuthorID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "author_id is required")
	}
	return nil
}

func (h *BlogHandler) validateUpdateBlogRequest(req service.UpdateBlogRequest) error {
	if req.Title != "" && (len(req.Title) < 2 || len(req.Title) > 200) {
		return echo.NewHTTPError(http.StatusBadRequest, "title must be between 2 and 200 characters")
	}
	if req.Content != "" && len(req.Content) < 10 {
		return echo.NewHTTPError(http.StatusBadRequest, "content must be at least 10 characters")
	}
	return nil
}

// SetupRoutes configures all API routes for blogs
func (h *BlogHandler) SetupRoutes(api *echo.Group) {
	v1 := api.Group("/v1")
	h.setupV1Routes(v1)
}

// setupV1Routes configures v1 API routes for blogs
func (h *BlogHandler) setupV1Routes(v1 *echo.Group) {
	blogs := v1.Group("/blogs")
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
