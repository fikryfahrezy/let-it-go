package http_server

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

// CustomValidator implements echo.Validator interface
type CustomValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

// NewCustomValidator creates a new CustomValidator instance
func NewCustomValidator() *CustomValidator {
	validate := validator.New()

	// Register field name function to use json tags as field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Setup translator
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		slog.Error("Failed to register translations", slog.String("error", err.Error()))
	}

	return &CustomValidator{
		validator: validate,
		trans:     trans,
	}
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i any) error {
	err := cv.validator.Struct(i)
	if err == nil {
		return nil
	}

	errorFields := make(map[string]any)

	for _, err := range err.(validator.ValidationErrors) {
		// Use Namespace() for nested field support including arrays
		// Examples: "User.Address.Street", "User.Addresses[0].Street", "Users[1].Name"
		fieldName := err.Field()
		if err.Namespace() != err.Field() {
			// For nested fields (including arrays), use the namespace but remove the root struct name
			namespace := err.Namespace()
			if dotIndex := strings.Index(namespace, "."); dotIndex != -1 {
				fieldName = namespace[dotIndex+1:]
			}
		}

		message := err.Translate(cv.trans)

		// Always use array format for consistency
		if existing, exists := errorFields[fieldName]; exists {
			errorFields[fieldName] = append(existing.([]string), message)
		} else {
			errorFields[fieldName] = []string{message}
		}
	}

	return &ValidationError{Fields: errorFields}
}

// ValidationError represents validation errors
type ValidationError struct {
	Fields map[string]any
}

func (ve *ValidationError) Error() string {
	return "validation failed"
}

// HandleValidationError checks if an error is a ValidationError and returns appropriate response
func HandleValidationError(c echo.Context, err error) error {
	if ve, ok := err.(*ValidationError); ok {
		return ValidationErrorResponse(c, "Validation failed", ve.Fields)
	}
	return err
}
