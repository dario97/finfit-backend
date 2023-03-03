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
	fieldValidationErrorMessage  = "some fields are invalid"
	bodyIsInvalidErrorMessage    = "body is invalid"
	paramsAreInvalidErrorMessage = "params are invalid, query params start_date and end_date are required"
	unexpectedErrorMessage       = "unexpected error"

	dateFormat = "2006-01-02"
)

type Handler interface {
	Add(context echo.Context) error
	SearchInPeriod(ctx echo.Context) error
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
		return h.buildErrorResponse(context, http.StatusBadRequest, bodyIsInvalidErrorMessage, err.Error())
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestBody); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrors)
	}

	command, err := h.mapAddCommandFromRequestBody(*requestBody)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, err.Error())
	}

	createdExpense, err := h.service.Add(command)
	if err != nil {
		return h.manageServiceError(context, err)
	}

	return context.JSON(http.StatusCreated, h.mapCreatedExpenseToExpenseResponse(createdExpense))
}

func (h handler) SearchInPeriod(context echo.Context) error {
	requestParams := new(searchInPeriodQueryParams)

	if err := context.Bind(requestParams); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, paramsAreInvalidErrorMessage, err.Error())
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestParams); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrors)
	}

	command, err := h.mapSearchCommandFromRequestBody(*requestParams)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, err.Error())
	}

	expenses, err := h.service.SearchInPeriod(command)
	if err != nil {
		return h.manageServiceError(context, err)
	}

	return context.JSON(http.StatusOK, h.mapExpensesToSearchResponse(expenses))

}

func (h handler) mapAddCommandFromRequestBody(body addExpenseRequest) (*expense.AddCommand, error) {
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

func (h handler) mapSearchCommandFromRequestBody(params searchInPeriodQueryParams) (*expense.SearchInPeriodCommand, error) {
	startDate, _ := time.Parse(dateFormat, params.StartDate)
	endDate, _ := time.Parse(dateFormat, params.EndDate)

	return expense.NewSearchInPeriodCommand(startDate, endDate)
}

func (h handler) mapCreatedExpenseToExpenseResponse(expense *models.Expense) expenseResponse {
	return expenseResponse{Expense: expenseBody{
		ID:          expense.Id.String(),
		Amount:      expense.Amount,
		ExpenseDate: expense.ExpenseDate.Format(dateFormat),
		Description: expense.Description,
		ExpenseType: expenseTypeBody{
			ID:   expense.ExpenseType.Id.String(),
			Name: expense.ExpenseType.Name,
		},
	}}
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

func (h handler) mapExpensesToSearchResponse(expenses []*models.Expense) searchResponse {
	expenseBodies := []expenseBody{}
	for _, expense := range expenses {
		expenseBodies = append(expenseBodies, h.mapExpenseToExpenseBody(expense))
	}

	return searchResponse{Expenses: expenseBodies}
}

func (h handler) mapExpenseToExpenseBody(expense *models.Expense) expenseBody {
	return expenseBody{
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

type addExpenseRequest struct {
	Amount      float64          `json:"amount,omitempty" validate:"required,gt=0"`
	ExpenseDate string           `json:"expense_date,omitempty" validate:"required,datetime=2006-01-02"`
	Description string           `json:"description,omitempty"`
	ExpenseType *expenseTypeBody `json:"expense_type,omitempty" validate:"required"`
}

type searchInPeriodQueryParams struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02,lteStrDateField=EndDate0x2C2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

type expenseTypeBody struct {
	ID   string `json:"id" validate:"required,uuid"`
	Name string `json:"name" validate:"required"`
}

type expenseBody struct {
	ID          string          `json:"id"`
	Amount      float64         `json:"amount"`
	ExpenseDate string          `json:"expense_date"`
	Description string          `json:"description"`
	ExpenseType expenseTypeBody `json:"expense_type"`
}

type expenseResponse struct {
	Expense expenseBody `json:"expense"`
}

type searchResponse struct {
	Expenses []expenseBody `json:"expenses"`
}
