package interfaces

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()

	return &CustomValidator{
		Validator: v,
	}
}

// Validate : Validate Data
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return translateValidationErrors(validationErrs)
		}
		return err
	}
	return nil
}

func translateValidationErrors(validationErrs validator.ValidationErrors) error {
	var messages []string
	for _, e := range validationErrs {
		fieldName := strcase.SnakeCase(e.Field())
		tag := e.Tag()
		param := e.Param()
		if len(param) > 0 {
			tag = fmt.Sprintf("%s=%s", tag, param)
		}
		value := e.Value()
		message := fmt.Sprintf(
			"invalid input on field '%s'; expected '%s', got '%v'",
			fieldName, tag, value,
		)
		messages = append(messages, message)
	}

	return errors.New(strings.Join(messages, ", "))
}
