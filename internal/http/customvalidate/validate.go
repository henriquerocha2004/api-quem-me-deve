package customvalidate

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/validateErrors"
)

var v *validator.Validate

type ValidationResponse struct {
	Errors []validateErrors.ValidationError `json:"errors"`
}

func init() {
	v = validator.New()
	_ = v.RegisterValidation("dateFormat:YYYY-MM-DD", dateValidate)
	_ = v.RegisterValidation("dateTimeFormat:YYYY-MM-DDTHH:MM:SS", dateTimeValidate)
}

func Validate(data any) *ValidationResponse {
	var errors []validateErrors.ValidationError
	validationResponse := &ValidationResponse{
		Errors: errors,
	}

	err := v.Struct(data)
	if err == nil {
		return validationResponse
	}

	validationErrors := err.(validator.ValidationErrors)
	for _, e := range validationErrors {
		errors = append(errors, validateErrors.ValidationError{
			Field:   e.Field(),
			Message: getErrorMessage(e),
		})
	}

	validationResponse.Errors = errors
	return validationResponse
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "dateFormat:YYYY-MM-DD":
		return "Invalid date format, expected YYYY-MM-DD"
	case "dateTimeFormat:YYYY-MM-DDTHH:MM:SS":
		return "Invalid date time format, expected YYYY-MM-DDTHH:MM:SS"
	default:
		return "Invalid value"
	}
}

func dateValidate(field validator.FieldLevel) bool {
	date := field.Field().String()
	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}

func dateTimeValidate(field validator.FieldLevel) bool {
	date := field.Field().String()
	_, err := time.Parse(time.DateTime, date)
	return err == nil
}
