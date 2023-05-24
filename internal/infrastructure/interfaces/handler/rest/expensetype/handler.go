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
	FieldValidationErrorMessage = "some fields are invalid"
	BodyIsInvalidErrorMessage   = "body is invalid"
	UnexpectedErrorMessage      = "unexpected error"
)

type Handler interface {
	Add(context echo.Context) error
	GetAll(context echo.Context) error
}

type handler struct {
	service         expensetype.Service
	fieldsValidator fieldvalidation.FieldsValidator
}

func NewHandler(service expensetype.Service, fieldsValidator fieldvalidation.FieldsValidator) *handler {
	return &handler{service: service, fieldsValidator: fieldsValidator}
}

func (h handler) Add(context echo.Context) error {
	requestBody := new(AddExpenseTypeRequest)

	if err := context.Bind(requestBody); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, BodyIsInvalidErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestBody); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, FieldValidationErrorMessage, FieldValidationErrorMessage, fieldValidationErrors, rest.FieldValidationErrorCode)
	}

	command, err := h.mapAddCommandFromRequestBody(*requestBody)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, FieldValidationErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	addedExpenseType, err := h.service.Add(command)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusInternalServerError, UnexpectedErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	return context.JSON(http.StatusCreated, h.mapAddedExpenseTypeToExpenseTypeResponse(addedExpenseType))
}

func (h handler) GetAll(context echo.Context) error {
	expenseTypes, err := h.service.GetAll()
	if err != nil {
		return h.buildErrorResponse(context, http.StatusInternalServerError, UnexpectedErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	return context.JSON(http.StatusOK, h.mapExpenseTypesToGetAllResponse(expenseTypes))
}

func (h handler) mapAddCommandFromRequestBody(body AddExpenseTypeRequest) (*expensetype.AddCommand, error) {
	return expensetype.NewAddCommand(body.Name)
}

func (h handler) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail string, fieldErrors []fieldvalidation.FieldError, errorCode uint) error {
	errorResponse := rest.ErrorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail, FieldErrors: fieldErrors, ErrorCode: errorCode}
	return ctx.JSON(statusCode, errorResponse)
}

func (h handler) mapAddedExpenseTypeToExpenseTypeResponse(expenseType *models.ExpenseType) AddExpenseTypeResponse {
	return AddExpenseTypeResponse{
		ExpenseType: h.mapExpenseTypeToExpenseTypeBody(expenseType),
	}
}

func (h handler) mapExpenseTypesToGetAllResponse(expenseTypes []*models.ExpenseType) GetAllResponse {
	expenseTypeBodies := []Body{}
	for _, expenseType := range expenseTypes {
		expenseTypeBodies = append(expenseTypeBodies, h.mapExpenseTypeToExpenseTypeBody(expenseType))
	}

	return GetAllResponse{ExpenseTypes: expenseTypeBodies}
}

func (h handler) mapExpenseTypeToExpenseTypeBody(expenseType *models.ExpenseType) Body {
	return Body{
		ID:   expenseType.Id().String(),
		Name: expenseType.Name(),
	}
}

type AddExpenseTypeRequest struct {
	Name string `json:"name,omitempty" validate:"required,min=3,max=32"`
}

type AddExpenseTypeResponse struct {
	ExpenseType Body `json:"expense_type"`
}

type GetAllResponse struct {
	ExpenseTypes []Body `json:"expense_types"`
}

type Body struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
