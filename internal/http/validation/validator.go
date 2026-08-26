package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"go-boilerplate/internal/http/response"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	english := en.New()
	uni := ut.New(english, english)
	translator, _ = uni.GetTranslator("en")

	validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return field.Name
		}

		return name
	})
	_ = enTranslations.RegisterDefaultTranslations(validate, translator)
}

func Validate(value any) []response.FieldError {
	err := validate.Struct(value)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []response.FieldError{{Field: "body", Message: err.Error()}}
	}

	fieldErrors := make([]response.FieldError, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		fieldErrors = append(fieldErrors, response.FieldError{
			Field:   fieldErr.Field(),
			Message: fieldErr.Translate(translator),
		})
	}

	return fieldErrors
}
