package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListAPIResponse represents a paginated API response
type ListAPIResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
	Error      string `json:"error,omitempty"`
}

// SwaggerErrorResponse represents an error response for Swagger docs
type SwaggerErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Legacy aliases for backward compatibility
type Response = APIResponse
type ListResponse = ListAPIResponse

func SuccessResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ListSuccessResponse(c echo.Context, message string, data any, pagination any) error {
	return c.JSON(http.StatusOK, ListResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

func BadRequestResponse(c echo.Context, message string, err error) error {
	return ErrorResponse(c, http.StatusBadRequest, message, err)
}

func NotFoundResponse(c echo.Context, message string, err error) error {
	return ErrorResponse(c, http.StatusNotFound, message, err)
}

func InternalServerErrorResponse(c echo.Context, message string, err error) error {
	return ErrorResponse(c, http.StatusInternalServerError, message, err)
}
