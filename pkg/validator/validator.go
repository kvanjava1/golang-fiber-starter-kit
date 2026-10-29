package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"golang-starter-kit/pkg/response"
)

type CustomValidator struct {
	validate *validator.Validate
}

var defaultValidator = New()

func New() *CustomValidator {
	v := validator.New()

	// Menggunakan tag JSON sebagai nama field jika tersedia
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	return &CustomValidator{validate: v}
}

// Validate memvalidasi struct dan mengembalikan slice of response.FieldError jika gagal
func (cv *CustomValidator) Validate(i any) []response.FieldError {
	err := cv.validate.Struct(i)
	if err == nil {
		return nil
	}

	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []response.FieldError{
			{Field: "general", Message: err.Error()},
		}
	}

	var fieldErrors []response.FieldError
	for _, fe := range valErrors {
		fieldErrors = append(fieldErrors, response.FieldError{
			Field:   fe.Field(),
			Message: formatMessage(fe),
		})
	}

	return fieldErrors
}

// ValidateStruct shortcut global menggunakan validator default
func ValidateStruct(i any) []response.FieldError {
	return defaultValidator.Validate(i)
}

func formatMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters or value", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters or value", fe.Field(), fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	default:
		return fmt.Sprintf("%s failed on '%s' validation", fe.Field(), fe.Tag())
	}
}
