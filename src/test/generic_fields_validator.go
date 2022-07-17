package test

import (
	validators2 "finfit-backend/src/infrastructure/interfaces/controller/validators"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type genericFieldsValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func newGenericFieldsValidator(validator *validator.Validate, translator ut.Translator) validators2.FieldsValidator {
	return genericFieldsValidator{
		validator:  validator,
		translator: translator,
	}
}

func (receiver genericFieldsValidator) ValidateFields(s interface{}) []validators2.FieldValidationError {
	var fieldValidationErrors []validators2.FieldValidationError
	if err := receiver.validator.Struct(s); err != nil {
		fieldValidationErrors = receiver.buildFieldValidationErrors(err.(validator.ValidationErrors))
	}

	return fieldValidationErrors
}

func (receiver genericFieldsValidator) buildFieldValidationErrors(fieldErrors []validator.FieldError) []validators2.FieldValidationError {
	var fieldValidationErrors []validators2.FieldValidationError
	for _, validationError := range fieldErrors {
		fieldValidationError := validators2.FieldValidationError{Message: validationError.Translate(receiver.translator),
			Field: validationError.Namespace()}
		fieldValidationErrors = append(fieldValidationErrors, fieldValidationError)
	}

	return fieldValidationErrors
}
