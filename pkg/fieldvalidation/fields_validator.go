package fieldvalidation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type FieldsValidator interface {
	ValidateFields(s interface{}) []FieldError
}

type fieldsValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func RegisterFieldsValidator() *fieldsValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	RegisterValidations(validate)
	_ = RegisterTranslations(validate, translator)

	return &fieldsValidator{
		validator:  validate,
		translator: translator,
	}
}

func (receiver fieldsValidator) ValidateFields(s interface{}) []FieldError {

	var fieldValidationErrors []FieldError
	if err := receiver.validator.Struct(s); err != nil {
		fieldValidationErrors = receiver.buildFieldValidationErrors(err.(validator.ValidationErrors))
	}

	return fieldValidationErrors
}

func (receiver fieldsValidator) buildFieldValidationErrors(fieldErrors []validator.FieldError) []FieldError {
	var fieldValidationErrors []FieldError
	for _, validationError := range fieldErrors {
		fieldValidationError := FieldError{Message: validationError.Translate(receiver.translator),
			Field: validationError.Field()}
		fieldValidationErrors = append(fieldValidationErrors, fieldValidationError)
	}

	return fieldValidationErrors
}
