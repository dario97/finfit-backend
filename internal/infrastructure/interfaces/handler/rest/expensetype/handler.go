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
	requestBody := new(addExpenseTypeRequest)

	if err := context.Bind(requestBody); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, bodyIsInvalidErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestBody); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrorMessage, fieldValidationErrors, rest.FieldValidationErrorCode)
	}

	command, err := h.mapAddCommandFromRequestBody(*requestBody)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	addedExpenseType, err := h.service.Add(command)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusInternalServerError, unexpectedErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	return context.JSON(http.StatusCreated, h.mapAddedExpenseTypeToExpenseTypeResponse(addedExpenseType))
}

func (h handler) GetAll(context echo.Context) error {
	expenseTypes, err := h.service.GetAll()
	if err != nil {
		return h.buildErrorResponse(context, http.StatusInternalServerError, unexpectedErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	return context.JSON(http.StatusOK, h.mapExpenseTypesToGetAllResponse(expenseTypes))
}

func (h handler) mapAddCommandFromRequestBody(body addExpenseTypeRequest) (*expensetype.AddCommand, error) {
	return expensetype.NewAddCommand(body.Name)
}

func (h handler) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail string, fieldErrors []fieldvalidation.FieldError, errorCode uint) error {
	errorResponse := rest.ErrorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail, FieldErrors: fieldErrors, ErrorCode: errorCode}
	return ctx.JSON(statusCode, errorResponse)
}

func (h handler) mapAddedExpenseTypeToExpenseTypeResponse(expenseType *models.ExpenseType) addExpenseTypeResponse {
	return addExpenseTypeResponse{
		ExpenseType: h.mapExpenseTypeToExpenseTypeBody(expenseType),
	}
}

func (h handler) mapExpenseTypesToGetAllResponse(expenseTypes []*models.ExpenseType) getAllResponse {
	expenseTypeBodies := []expenseTypeBody{}
	for _, expenseType := range expenseTypes {
		expenseTypeBodies = append(expenseTypeBodies, h.mapExpenseTypeToExpenseTypeBody(expenseType))
	}

	return getAllResponse{ExpenseTypes: expenseTypeBodies}
}

func (h handler) mapExpenseTypeToExpenseTypeBody(expenseType *models.ExpenseType) expenseTypeBody {
	return expenseTypeBody{
		ID:   expenseType.Id.String(),
		Name: expenseType.Name,
	}
}

type addExpenseTypeRequest struct {
	Name string `json:"name,omitempty" validate:"required,min=3,max=32"`
}

type addExpenseTypeResponse struct {
	ExpenseType expenseTypeBody `json:"expense_type"`
}

type getAllResponse struct {
	ExpenseTypes []expenseTypeBody `json:"expense_types"`
}

type expenseTypeBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
