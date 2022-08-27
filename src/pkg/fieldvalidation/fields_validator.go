package fieldvalidation

type FieldsValidator interface {
	ValidateFields(s interface{}) []FieldValidationError
}
