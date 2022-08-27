package searchinperiod

import (
	"finfit-backend/src/pkg/fieldvalidation"
	"github.com/labstack/echo"
)

type Handler interface {
	SearchInPeriod(context echo.Context) error
}

type handler struct {
	service         Service
	fieldsValidator fieldvalidation.FieldsValidator
}

func NewHandler(service Service, fieldsValidator fieldvalidation.FieldsValidator) Handler {
	return handler{
		service:         service,
		fieldsValidator: fieldsValidator,
	}
}

func (h handler) SearchInPeriod(context echo.Context) error {
	//TODO implement me
	panic("implement me")
}
