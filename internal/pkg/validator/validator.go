package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/syst3mctl/check-in-api/internal/core/domain"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
	// Register function to get json tag name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s interface{}) *domain.ErrorResponse {
	if validate == nil {
		Init()
	}

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string][]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		var msg string
		switch err.Tag() {
		case "required":
			msg = "field is required"
		case "email":
			msg = "email is invalid format"
		case "min":
			msg = fmt.Sprintf("must be at least %s characters", err.Param())
		case "max":
			msg = fmt.Sprintf("must be at most %s characters", err.Param())
		case "oneof":
			msg = fmt.Sprintf("must be one of: %s", strings.ReplaceAll(err.Param(), " ", ", "))
		default:
			msg = fmt.Sprintf("failed on tag %s", err.Tag())
		}
		validationErrors[field] = append(validationErrors[field], msg)
	}

	return &domain.ErrorResponse{
		Message: "invalid payload",
		Errors:  validationErrors,
	}
}
