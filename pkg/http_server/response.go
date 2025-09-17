package server

import (
	"net/http"

	"github.com/fikryfahrezy/let-it-go/pkg/app_error"
	"github.com/labstack/echo/v4"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

// ListAPIResponse represents a paginated API response
type ListAPIResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
	Error      string `json:"error,omitempty"`
	ErrorCode  string `json:"error_code,omitempty"`
}

func SuccessResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ListSuccessResponse(c echo.Context, message string, data any, pagination any) error {
	return c.JSON(http.StatusOK, ListAPIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string, err error) error {
	errorMsg := ""
	errorCode := ""
	
	if err != nil {
		errorMsg = err.Error()
		// Extract error code if it's an AppError
		errorCode = app_error.GetCode(err)
	}

	return c.JSON(statusCode, APIResponse{
		Success:   false,
		Message:   message,
		Error:     errorMsg,
		ErrorCode: errorCode,
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