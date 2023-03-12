package rest

import "finfit-backend/pkg/fieldvalidation"

const (
	FieldValidationErrorCode = 1
)

type ErrorResponse struct {
	StatusCode  int                          `json:"status_code"`
	Msg         string                       `json:"msg"`
	ErrorDetail string                       `json:"error_detail"`
	FieldErrors []fieldvalidation.FieldError `json:"field_errors"`
	ErrorCode   uint                         `json:"error_code"`
}
