package search_in_period

import (
	"finfit-backend/src/pkg/validator"
	"github.com/labstack/echo"
)

type Handler interface {
	SearchInPeriod(context echo.Context) error
}

type handler struct {
	service         Service
	fieldsValidator validator.FieldsValidator
}

func NewHandler(service Service, fieldsValidator validator.FieldsValidator) Handler {
	return handler{
		service:         service,
		fieldsValidator: fieldsValidator,
	}
}

func (h handler) SearchInPeriod(context echo.Context) error {
	//TODO implement me
	panic("implement me")
}
