package validator

type FieldsValidator interface {
	ValidateFields(s interface{}) []FieldValidationError
}
