package expense

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expense"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest"
	"finfit-backend/pkg/fieldvalidation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	fieldValidationErrorMessage = "some fields are invalid"
	bodyWasInvalidErrorMessage  = "body was invalid"
	unexpectedErrorMessage      = "unexpected error"
	dateFormat                  = "02-01-2006"
)

type Handler interface {
	Add(context echo.Context) error
}

type handler struct {
	service         expense.Service
	fieldsValidator fieldvalidation.FieldsValidator
}

func NewHandler(service expense.Service, fieldsValidator fieldvalidation.FieldsValidator) Handler {
	return handler{
		service:         service,
		fieldsValidator: fieldsValidator,
	}
}

func (h handler) Add(context echo.Context) error {
	requestBody := new(addExpenseRequest)

	if err := context.Bind(requestBody); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, bodyWasInvalidErrorMessage, err.Error())
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestBody); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrors)
	}

	command, err := h.mapCommandFromRequestBody(*requestBody)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, err.Error())
	}

	createdExpense, err := h.service.Add(command)
	if err != nil {
		return h.manageServiceError(context, err)
	}

	return context.JSON(http.StatusCreated, h.mapCreatedExpenseToExpenseResponse(createdExpense))
}

func (h handler) mapCommandFromRequestBody(body addExpenseRequest) (*expense.AddCommand, error) {
	date, _ := time.Parse(dateFormat, body.ExpenseDate)
	expenseTypeId, err := uuid.Parse(body.ExpenseType.ID)
	if err != nil {
		return nil, err
	}

	return expense.NewAddCommand(body.Amount, date, body.Description, &models.ExpenseType{
		Id:   expenseTypeId,
		Name: body.ExpenseType.Name,
	})
}

func (h handler) mapCreatedExpenseToExpenseResponse(expense *models.Expense) expenseResponse {
	return expenseResponse{
		ID:          expense.Id.String(),
		Amount:      expense.Amount,
		ExpenseDate: expense.ExpenseDate.Format(dateFormat),
		Description: expense.Description,
		ExpenseType: expenseTypeBody{
			ID:   expense.ExpenseType.Id.String(),
			Name: expense.ExpenseType.Name,
		},
	}
}

func (h handler) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail interface{}) error {
	errorResponse := rest.ErrorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail}
	return ctx.JSON(statusCode, errorResponse)
}

func (h handler) manageServiceError(ctx echo.Context, err error) error {
	if errors.As(err, &expense.InvalidExpenseTypeError{}) {
		return h.buildErrorResponse(ctx, http.StatusBadRequest, err.Error(), err.Error())
	} else {
		return h.buildErrorResponse(ctx, http.StatusInternalServerError, unexpectedErrorMessage, err.Error())
	}
}

type addExpenseRequest struct {
	Amount      float64          `json:"amount" validate:"required,gt=0"`
	ExpenseDate string           `json:"expense_date" validate:"required,datetime=02-01-2006"`
	Description string           `json:"description"`
	ExpenseType *expenseTypeBody `json:"expense_type" validate:"required"`
}

type expenseTypeBody struct {
	ID   string `json:"id" validate:"required,uuid"`
	Name string `json:"name" validate:"required"`
}

type expenseResponse struct {
	ID          string          `json:"id"`
	Amount      float64         `json:"amount"`
	ExpenseDate string          `json:"expense_date"`
	Description string          `json:"description"`
	ExpenseType expenseTypeBody `json:"expense_type"`
}
