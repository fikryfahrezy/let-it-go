package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fikryfahrezy/let-it-go/feature/user/handler"
	"github.com/fikryfahrezy/let-it-go/feature/user/repository"
	"github.com/fikryfahrezy/let-it-go/feature/user/service"
	"github.com/fikryfahrezy/let-it-go/feature/user/service/servicefakes"
	"github.com/fikryfahrezy/let-it-go/pkg/http_server"
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

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userID := uuid.New()
	expectedResponse := service.CreateUserResponse{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockService.CreateUserReturns(expectedResponse, nil)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	requestBody := service.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userHandler.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.CreateUserCallCount())
	_, actualReq := mockService.CreateUserArgsForCall(0)
	assert.Equal(t, requestBody.Name, actualReq.Name)
	assert.Equal(t, requestBody.Email, actualReq.Email)
	assert.Equal(t, requestBody.Password, actualReq.Password)
}

func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	// Missing required fields
	requestBody := service.CreateUserRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userHandler.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Service should not be called on validation error
	assert.Equal(t, 0, mockService.CreateUserCallCount())
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	mockService.CreateUserReturns(service.CreateUserResponse{}, service.ErrUserAlreadyExists)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	requestBody := service.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userHandler.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.CreateUserCallCount())
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userID := uuid.New()
	expectedResponse := service.GetUserResponse{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockService.GetUserByIDReturns(expectedResponse, nil)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	err := userHandler.GetUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with correct ID
	assert.Equal(t, 1, mockService.GetUserByIDCallCount())
	_, actualID := mockService.GetUserByIDArgsForCall(0)
	assert.Equal(t, userID, actualID)
}

func TestUserHandler_GetUser_InvalidUUID(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid-uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	err := userHandler.GetUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Service should not be called on invalid UUID
	assert.Equal(t, 0, mockService.GetUserByIDCallCount())
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userID := uuid.New()
	mockService.GetUserByIDReturns(service.GetUserResponse{}, repository.ErrUserNotFound)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	err := userHandler.GetUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.GetUserByIDCallCount())
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	expectedUsers := []service.ListUsersResponse{
		{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Jane Doe", 
			Email:     "jane@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.ListUsersReturns(expectedUsers, 2, nil)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.ListUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with default pagination
	assert.Equal(t, 1, mockService.ListUsersCallCount())
	_, paginationReq := mockService.ListUsersArgsForCall(0)
	assert.Equal(t, 1, paginationReq.Page)
	assert.Equal(t, 10, paginationReq.PageSize)
}

func TestUserHandler_ListUsers_WithPagination(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	expectedUsers := []service.ListUsersResponse{
		{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.ListUsersReturns(expectedUsers, 1, nil)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=2&page_size=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.ListUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with custom pagination
	assert.Equal(t, 1, mockService.ListUsersCallCount())
	_, paginationReq := mockService.ListUsersArgsForCall(0)
	assert.Equal(t, 2, paginationReq.Page)
	assert.Equal(t, 5, paginationReq.PageSize)
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userID := uuid.New()
	mockService.DeleteUserReturns(nil)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	err := userHandler.DeleteUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify service was called with correct ID
	assert.Equal(t, 1, mockService.DeleteUserCallCount())
	_, actualID := mockService.DeleteUserArgsForCall(0)
	assert.Equal(t, userID, actualID)
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userID := uuid.New()
	mockService.DeleteUserReturns(repository.ErrUserNotFound)

	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	err := userHandler.DeleteUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Verify service was called
	assert.Equal(t, 1, mockService.DeleteUserCallCount())
}

func TestUserHandler_HealthCheck(t *testing.T) {
	mockService := &servicefakes.FakeUserService{}
	userHandler := handler.NewUserHandler(mockService)
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := userHandler.HealthCheck(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Service is healthy", response["message"])

	// Health check should not call service methods
	assert.Equal(t, 0, mockService.CreateUserCallCount())
	assert.Equal(t, 0, mockService.GetUserByIDCallCount())
	assert.Equal(t, 0, mockService.UpdateUserCallCount())
	assert.Equal(t, 0, mockService.DeleteUserCallCount())
	assert.Equal(t, 0, mockService.ListUsersCallCount())
}