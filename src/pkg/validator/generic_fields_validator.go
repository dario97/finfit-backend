package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type genericFieldsValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewGenericFieldsValidator(validator *validator.Validate, translator ut.Translator) FieldsValidator {
	return genericFieldsValidator{
		validator:  validator,
		translator: translator,
	}
}

func (receiver genericFieldsValidator) ValidateFields(s interface{}) []FieldValidationError {
	var fieldValidationErrors []FieldValidationError
	if err := receiver.validator.Struct(s); err != nil {
		fieldValidationErrors = receiver.buildFieldValidationErrors(err.(validator.ValidationErrors))
	}

	return fieldValidationErrors
}

func (receiver genericFieldsValidator) buildFieldValidationErrors(fieldErrors []validator.FieldError) []FieldValidationError {
	var fieldValidationErrors []FieldValidationError
	for _, validationError := range fieldErrors {
		fieldValidationError := FieldValidationError{Message: validationError.Translate(receiver.translator),
			Field: validationError.Field()}
		fieldValidationErrors = append(fieldValidationErrors, fieldValidationError)
	}

	return fieldValidationErrors
}
