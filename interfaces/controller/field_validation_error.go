package controller

type FieldValidationError struct {
	Field            string `json:"field"`
	ValidationResult string `json:"validation_result"`
}
