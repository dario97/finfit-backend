package test

import (
	"finfit-backend/src/domain/validators"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type genericFieldsValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func newGenericFieldsValidator(validator *validator.Validate, translator ut.Translator) genericFieldsValidator {
	return genericFieldsValidator{
		validator:  validator,
		translator: translator,
	}
}

func (receiver genericFieldsValidator) ValidateFields(s interface{}) []validators.FieldValidationError {
	var fieldValidationErrors []validators.FieldValidationError
	if err := receiver.validator.Struct(s); err != nil {
		fieldValidationErrors = receiver.buildFieldValidationErrors(err.(validator.ValidationErrors))
	}

	return fieldValidationErrors
}

func (receiver genericFieldsValidator) buildFieldValidationErrors(fieldErrors []validator.FieldError) []validators.FieldValidationError {
	var fieldValidationErrors []validators.FieldValidationError
	for _, validationError := range fieldErrors {
		fieldValidationError := validators.FieldValidationError{Message: validationError.Translate(receiver.translator),
			Field: validationError.Namespace()}
		fieldValidationErrors = append(fieldValidationErrors, fieldValidationError)
	}

	return fieldValidationErrors
}
