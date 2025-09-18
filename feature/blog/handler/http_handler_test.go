package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/blog/handler"
	"github.com/fikryfahrezy/let-it-go/feature/blog/repository"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service"
	"github.com/fikryfahrezy/let-it-go/feature/blog/service/servicefakes"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
	"github.com/fikryfahrezy/let-it-go/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = http_server.NewCustomValidator()
	return e
}

func TestBlogHandler_CreateBlog_Success(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogID := uuid.New()
	authorID := uuid.New()
	expectedResponse := service.GetBlogResponse{
		ID:        blogID,
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		AuthorID:  authorID,
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockService.CreateBlogReturns(expectedResponse, nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	requestBody := service.CreateBlogRequest{
		Title:    "Test Blog",
		Content:  "This is a test blog content",
		AuthorID: authorID,
		Status:   "draft",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/blogs", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = blogHandler.CreateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.CreateBlogCallCount())
	_, actualReq := mockService.CreateBlogArgsForCall(0)
	assert.Equal(t, requestBody.Title, actualReq.Title)
	assert.Equal(t, requestBody.Content, actualReq.Content)
	assert.Equal(t, requestBody.AuthorID, actualReq.AuthorID)
	assert.Equal(t, requestBody.Status, actualReq.Status)
}

func TestBlogHandler_CreateBlog_ValidationError(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	// Missing required fields
	requestBody := service.CreateBlogRequest{
		Title:   "Test Blog",
		Content: "Short", // Too short content
		Status:  "draft",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/blogs", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = blogHandler.CreateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Service should not be called on validation error
	assert.Equal(t, 0, mockService.CreateBlogCallCount())
}

func TestBlogHandler_CreateBlog_InvalidStatus(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	authorID := uuid.New()
	requestBody := service.CreateBlogRequest{
		Title:    "Test Blog",
		Content:  "This is a test blog content",
		AuthorID: authorID,
		Status:   "invalid-status",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/blogs", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = blogHandler.CreateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Service should not be called on validation error
	assert.Equal(t, 0, mockService.CreateBlogCallCount())
}

func TestBlogHandler_GetBlog_Success(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogID := uuid.New()
	authorID := uuid.New()
	expectedResponse := service.GetBlogResponse{
		ID:        blogID,
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		AuthorID:  authorID,
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockService.GetBlogByIDReturns(expectedResponse, nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs/"+blogID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/:id")
	c.SetParamNames("id")
	c.SetParamValues(blogID.String())

	err := blogHandler.GetBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with correct ID
	assert.Equal(t, 1, mockService.GetBlogByIDCallCount())
	_, actualID := mockService.GetBlogByIDArgsForCall(0)
	assert.Equal(t, blogID, actualID)
}

func TestBlogHandler_GetBlog_InvalidUUID(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs/invalid-uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/:id")
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	err := blogHandler.GetBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Service should not be called on invalid UUID
	assert.Equal(t, 0, mockService.GetBlogByIDCallCount())
}

func TestBlogHandler_GetBlog_NotFound(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogID := uuid.New()
	mockService.GetBlogByIDReturns(service.GetBlogResponse{}, repository.ErrBlogNotFound)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs/"+blogID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/:id")
	c.SetParamNames("id")
	c.SetParamValues(blogID.String())

	err := blogHandler.GetBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.GetBlogByIDCallCount())
}

func TestBlogHandler_ListBlogs_Success(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	authorID := uuid.New()
	expectedBlogs := []service.GetBlogResponse{
		{
			ID:        uuid.New(),
			Title:     "Blog 1",
			Content:   "Content 1",
			AuthorID:  authorID,
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Blog 2",
			Content:   "Content 2",
			AuthorID:  authorID,
			Status:    "draft",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.ListBlogsReturns(expectedBlogs, 2, nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := blogHandler.ListBlogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with default pagination
	assert.Equal(t, 1, mockService.ListBlogsCallCount())
	_, paginationReq := mockService.ListBlogsArgsForCall(0)
	assert.Equal(t, 1, paginationReq.Page)
	assert.Equal(t, 10, paginationReq.PageSize)
}

func TestBlogHandler_ListBlogs_WithPagination(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	authorID := uuid.New()
	expectedBlogs := []service.GetBlogResponse{
		{
			ID:        uuid.New(),
			Title:     "Blog 1",
			Content:   "Content 1",
			AuthorID:  authorID,
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.ListBlogsReturns(expectedBlogs, 1, nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs?page=2&page_size=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := blogHandler.ListBlogs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with custom pagination
	assert.Equal(t, 1, mockService.ListBlogsCallCount())
	_, paginationReq := mockService.ListBlogsArgsForCall(0)
	assert.Equal(t, 2, paginationReq.Page)
	assert.Equal(t, 5, paginationReq.PageSize)
}

func TestBlogHandler_DeleteBlog_Success(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogID := uuid.New()
	mockService.DeleteBlogReturns(nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/blogs/"+blogID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/:id")
	c.SetParamNames("id")
	c.SetParamValues(blogID.String())

	err := blogHandler.DeleteBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with correct ID
	assert.Equal(t, 1, mockService.DeleteBlogCallCount())
	_, actualID := mockService.DeleteBlogArgsForCall(0)
	assert.Equal(t, blogID, actualID)
}

func TestBlogHandler_DeleteBlog_NotFound(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	blogID := uuid.New()
	mockService.DeleteBlogReturns(repository.ErrBlogNotFound)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/blogs/"+blogID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/:id")
	c.SetParamNames("id")
	c.SetParamValues(blogID.String())

	err := blogHandler.DeleteBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.DeleteBlogCallCount())
}

func TestBlogHandler_GetBlogsByAuthor_Success(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	authorID := uuid.New()
	expectedBlogs := []service.GetBlogResponse{
		{
			ID:        uuid.New(),
			Title:     "Author Blog 1",
			Content:   "Content 1",
			AuthorID:  authorID,
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Author Blog 2",
			Content:   "Content 2",
			AuthorID:  authorID,
			Status:    "draft",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.GetBlogsByAuthorReturns(expectedBlogs, 2, nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs/author/"+authorID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/author/:author_id")
	c.SetParamNames("author_id")
	c.SetParamValues(authorID.String())

	err := blogHandler.GetBlogsByAuthor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with correct author ID
	assert.Equal(t, 1, mockService.GetBlogsByAuthorCallCount())
	_, actualAuthorID, paginationReq := mockService.GetBlogsByAuthorArgsForCall(0)
	assert.Equal(t, authorID, actualAuthorID)
	assert.Equal(t, 1, paginationReq.Page)
	assert.Equal(t, 10, paginationReq.PageSize)
}

func TestBlogHandler_GetBlogsByStatus_Success(t *testing.T) {
	mockService := &servicefakes.FakeBlogService{}
	expectedBlogs := []service.GetBlogResponse{
		{
			ID:        uuid.New(),
			Title:     "Published Blog 1",
			Content:   "Content 1",
			AuthorID:  uuid.New(),
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Published Blog 2",
			Content:   "Content 2",
			AuthorID:  uuid.New(),
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.GetBlogsByStatusReturns(expectedBlogs, 2, nil)

	blogHandler := handler.NewBlogHandler(logger.NewDiscardLogger(), mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/blogs/status/published", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/blogs/status/:status")
	c.SetParamNames("status")
	c.SetParamValues("published")

	err := blogHandler.GetBlogsByStatus(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with correct status
	assert.Equal(t, 1, mockService.GetBlogsByStatusCallCount())
	_, actualStatus, paginationReq := mockService.GetBlogsByStatusArgsForCall(0)
	assert.Equal(t, "published", actualStatus)
	assert.Equal(t, 1, paginationReq.Page)
	assert.Equal(t, 10, paginationReq.PageSize)
}
