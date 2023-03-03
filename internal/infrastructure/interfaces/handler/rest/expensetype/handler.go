package expensetype

import (
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expensetype"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest"
	"finfit-backend/pkg/fieldvalidation"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	fieldValidationErrorMessage = "some fields are invalid"
	bodyIsInvalidErrorMessage   = "body is invalid"
	unexpectedErrorMessage      = "unexpected error"
)

type Handler interface {
	Add(context echo.Context) error
}

type handler struct {
	service         expensetype.Service
	fieldsValidator fieldvalidation.FieldsValidator
}

func NewHandler(service expensetype.Service, fieldsValidator fieldvalidation.FieldsValidator) *handler {
	return &handler{service: service, fieldsValidator: fieldsValidator}
}

func (h handler) Add(context echo.Context) error {
	requestBody := new(addExpenseTypeRequest)

	if err := context.Bind(requestBody); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, bodyIsInvalidErrorMessage, err.Error())
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestBody); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrors)
	}

	command, err := h.mapAddCommandFromRequestBody(*requestBody)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, err.Error())
	}

	addedExpenseType, err := h.service.Add(command)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusInternalServerError, unexpectedErrorMessage, err.Error())
	}

	return context.JSON(http.StatusCreated, h.mapAddedExpenseTypeToExpenseTypeResponse(addedExpenseType))
}

func (h handler) mapAddCommandFromRequestBody(body addExpenseTypeRequest) (*expensetype.AddCommand, error) {
	return expensetype.NewAddCommand(body.Name)
}

func (h handler) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail interface{}) error {
	errorResponse := rest.ErrorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail}
	return ctx.JSON(statusCode, errorResponse)
}

func (h handler) mapAddedExpenseTypeToExpenseTypeResponse(expenseType *models.ExpenseType) addExpenseTypeResponse {
	return addExpenseTypeResponse{
		ID:   expenseType.Id.String(),
		Name: expenseType.Name,
	}
}

type addExpenseTypeRequest struct {
	Name string `json:"name,omitempty" validate:"required,min=3,max=32"`
}

type addExpenseTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
