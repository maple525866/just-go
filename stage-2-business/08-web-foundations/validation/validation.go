package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"just-go/stage-2-business/08-web-foundations/model"
)

// Validator wraps go-playground/validator so handlers do not depend on its
// error types directly.
type Validator struct {
	validate *validator.Validate
}

// New builds the chapter validator.
func New() *Validator {
	return &Validator{validate: validator.New(validator.WithRequiredStructEnabled())}
}

// ValidateCreateArticle returns client-readable field errors.
func (v *Validator) ValidateCreateArticle(request model.CreateArticleRequest) []model.FieldError {
	if err := v.validate.Struct(request); err != nil {
		return fieldErrors(err)
	}
	return nil
}

func fieldErrors(err error) []model.FieldError {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []model.FieldError{{Field: "request", Rule: "invalid", Message: err.Error()}}
	}

	fields := make([]model.FieldError, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		field := jsonFieldName(validationError)
		fields = append(fields, model.FieldError{
			Field:   field,
			Rule:    validationError.Tag(),
			Message: fmt.Sprintf("%s failed %s validation", field, validationError.Tag()),
		})
	}
	return fields
}

func jsonFieldName(err validator.FieldError) string {
	field := err.Field()
	switch field {
	case "Title":
		field = "title"
	case "Body":
		field = "body"
	case "Tags":
		field = "tags"
	}

	namespace := err.StructNamespace()
	if strings.Contains(namespace, "Tags[") {
		start := strings.Index(namespace, "Tags[") + len("Tags")
		return "tags" + namespace[start:]
	}
	return field
}
