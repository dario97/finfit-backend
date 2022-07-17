package validators

type FieldsValidator interface {
	ValidateFields(s interface{}) []FieldValidationError
}
