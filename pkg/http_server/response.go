package http_server

import (
	"net/http"

	"github.com/fikryfahrezy/let-it-go/pkg/app_error"
	"github.com/labstack/echo/v4"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Message     string         `json:"message"`
	Error       string         `json:"error"`
	ErrorFields map[string]any `json:"error_fields"`
	Result      any            `json:"result"`
}

// ListAPIResponse represents a paginated API response
type ListAPIResponse struct {
	Message     string              `json:"message"`
	Error       string              `json:"error"`
	ErrorFields map[string]any      `json:"error_fields"`
	Result      any                 `json:"result"`
	Pagination  *PaginationResponse `json:"pagination,omitempty"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	TotalData  int `json:"total_data"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
}

// PaginationRequest represents pagination input parameters
type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func SuccessResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, APIResponse{
		Message: message,
		Error:   "",
		Result:  data,
	})
}

func CreatedResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Message: message,
		Error:   "",
		Result:  data,
	})
}

func ListSuccessResponse(c echo.Context, message string, data any, pagination PaginationResponse) error {
	return c.JSON(http.StatusOK, ListAPIResponse{
		Message:    message,
		Error:      "",
		Result:     data,
		Pagination: &pagination,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string, err error) error {
	errorCode := "UNKNOWN_ERROR"
	errorMessage := message
	var errorFields map[string]any

	if err != nil {
		// Extract error code if it's an AppError
		code := app_error.GetCode(err)
		if code != "" {
			errorCode = code
		}

		// Extract error message if it's an AppError
		message := app_error.GetMessage(err)
		if message != "" {
			errorMessage = message
		}

		// You can add error fields parsing here if needed
		errorFields = map[string]any{}
	}

	return c.JSON(statusCode, APIResponse{
		Message:     errorMessage,
		Error:       errorCode,
		ErrorFields: errorFields,
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

// ValidationErrorResponse creates a response for validation errors with field-specific errors
func ValidationErrorResponse(c echo.Context, message string, errorFields map[string]any) error {
	return c.JSON(http.StatusBadRequest, APIResponse{
		Message:     message,
		Error:       "VALIDATION_ERROR",
		ErrorFields: errorFields,
	})
}

// CreatePaginationResponse creates a pagination response object
func CreatePaginationResponse(totalData, totalPages, page, limit int) PaginationResponse {
	return PaginationResponse{
		TotalData:  totalData,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}
}
