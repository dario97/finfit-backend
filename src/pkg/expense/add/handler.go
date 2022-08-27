package add

import (
	"errors"
	"finfit-backend/src/pkg/fieldvalidation"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

const (
	fieldValidationErrorMessage = "some fields are invalid"
	bodyWasInvalidErrorMessage  = "body was invalid"
	unexpectedErrorMessage      = "unexpected error"
	dateFormat                  = "2006-01-02"
)

type Handler interface {
	Add(context echo.Context) error
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

func (h handler) Add(context echo.Context) error {
	requestBody := new(requestBody)

	if err := context.Bind(requestBody); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, bodyWasInvalidErrorMessage, err.Error())
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestBody); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrors)
	}

	createdExpense, err := h.service.Add(h.mapCommandFromRequestBody(*requestBody))

	if err != nil {
		return h.manageServiceError(context, err)
	}

	return context.JSON(http.StatusCreated, h.mapCreatedExpenseToExpenseResponse(*createdExpense))
}

func (h handler) mapCommandFromRequestBody(body requestBody) command {
	date, _ := time.Parse(dateFormat, body.ExpenseDate)
	return command{
		amount:      body.Amount,
		expenseDate: date,
		description: body.Description,
		expenseType: expenseType{
			id:   body.ExpenseType.ID,
			name: body.ExpenseType.Name,
		},
	}
}

func (h handler) mapCreatedExpenseToExpenseResponse(expense createdExpense) expenseResponse {
	return expenseResponse{
		ID:          expense.id,
		Amount:      expense.amount,
		ExpenseDate: expense.expenseDate.Format(dateFormat),
		Description: expense.description,
		ExpenseType: expenseTypeRequestBody{
			ID:   expense.expenseType.id,
			Name: expense.expenseType.name,
		},
	}
}

func (h handler) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail interface{}) error {
	errorResponse := errorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail}
	return ctx.JSON(statusCode, errorResponse)
}

func (h handler) manageServiceError(ctx echo.Context, err error) error {
	if errors.As(err, &InvalidExpenseTypeError{}) {
		return h.buildErrorResponse(ctx, http.StatusBadRequest, err.Error(), err.Error())
	} else {
		return h.buildErrorResponse(ctx, http.StatusInternalServerError, unexpectedErrorMessage, err.Error())
	}
}

type requestBody struct {
	Amount      float64                 `json:"amount" validate:"required,gt=0"`
	ExpenseDate string                  `json:"expense_date" validate:"required,datetime=2006-01-02"`
	Description string                  `json:"description"`
	ExpenseType *expenseTypeRequestBody `json:"expense_type" validate:"required"`
}

type expenseTypeRequestBody struct {
	ID   uint64 `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type expenseResponse struct {
	ID          uint64                 `json:"id"`
	Amount      float64                `json:"amount"`
	ExpenseDate string                 `json:"expense_date"`
	Description string                 `json:"description"`
	ExpenseType expenseTypeRequestBody `json:"expense_type"`
}

type errorResponse struct {
	StatusCode  int         `json:"status_code"`
	Msg         string      `json:"msg"`
	ErrorDetail interface{} `json:"error_detail"`
}
