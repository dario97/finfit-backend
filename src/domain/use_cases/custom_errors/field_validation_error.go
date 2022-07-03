package custom_errors

type FieldValidationError struct {
	Field            string `json:"field"`
	ValidationResult string `json:"validation_result"`
}
