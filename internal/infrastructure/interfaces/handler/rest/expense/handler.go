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
	FieldValidationErrorMessage  = "some fields are invalid"
	BodyIsInvalidErrorMessage    = "body is invalid"
	ParamsAreInvalidErrorMessage = "params are invalid, query params start_date and end_date are required"
	UnexpectedErrorMessage       = "unexpected error"
	DateFormat                   = "2006-01-02"
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
	requestBody := new(AddExpenseRequest)

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

	createdExpense, err := h.service.Add(command)
	if err != nil {
		return h.manageServiceError(context, err)
	}

	return context.JSON(http.StatusCreated, h.mapCreatedExpenseToExpenseResponse(createdExpense))
}

// TODO: Analizar limitar el largo del periodo y limite de expenses a buscar
func (h handler) SearchInPeriod(context echo.Context) error {
	requestParams := new(SearchInPeriodQueryParams)

	if err := context.Bind(requestParams); err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, ParamsAreInvalidErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	if fieldValidationErrors := h.fieldsValidator.ValidateFields(requestParams); len(fieldValidationErrors) > 0 {
		return h.buildErrorResponse(context, http.StatusBadRequest, FieldValidationErrorMessage, FieldValidationErrorMessage, fieldValidationErrors, rest.FieldValidationErrorCode)
	}

	command, err := h.mapSearchCommandFromRequestBody(*requestParams)
	if err != nil {
		return h.buildErrorResponse(context, http.StatusBadRequest, FieldValidationErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}

	expenses, err := h.service.SearchInPeriod(command)
	if err != nil {
		return h.manageServiceError(context, err)
	}

	return context.JSON(http.StatusOK, h.mapExpensesToSearchResponse(expenses))

}

func (h handler) mapAddCommandFromRequestBody(body AddExpenseRequest) (*expense.AddCommand, error) {
	date, _ := time.Parse(DateFormat, body.ExpenseDate)
	expenseTypeId, err := uuid.Parse(body.ExpenseType.ID)
	if err != nil {
		return nil, err
	}

	return expense.NewAddCommand(body.Amount.Amount, body.Amount.Currency, date, body.Description, expenseTypeId)
}

func (h handler) mapSearchCommandFromRequestBody(params SearchInPeriodQueryParams) (*expense.SearchInPeriodCommand, error) {
	startDate, _ := time.Parse(DateFormat, params.StartDate)
	endDate, _ := time.Parse(DateFormat, params.EndDate)

	return expense.NewSearchInPeriodCommand(startDate, endDate)
}

func (h handler) mapCreatedExpenseToExpenseResponse(expense *models.Expense) Response {
	return Response{Expense: h.mapExpenseToExpenseBody(expense)}
}

func (h handler) manageServiceError(ctx echo.Context, err error) error {
	if errors.As(err, &expense.InvalidExpenseTypeError{}) {
		return h.buildErrorResponse(ctx, http.StatusBadRequest, err.Error(), err.Error(), []fieldvalidation.FieldError{}, 0)
	} else {
		return h.buildErrorResponse(ctx, http.StatusInternalServerError, UnexpectedErrorMessage, err.Error(), []fieldvalidation.FieldError{}, 0)
	}
}

func (h handler) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail string, fieldErrors []fieldvalidation.FieldError, errorCode uint) error {
	errorResponse := rest.ErrorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail, FieldErrors: fieldErrors, ErrorCode: errorCode}
	return ctx.JSON(statusCode, errorResponse)
}

func (h handler) mapExpensesToSearchResponse(expenses []*models.Expense) SearchResponse {
	expenseBodies := []Body{}
	for _, expense := range expenses {
		expenseBodies = append(expenseBodies, h.mapExpenseToExpenseBody(expense))
	}

	return SearchResponse{Expenses: expenseBodies}
}

func (h handler) mapExpenseToExpenseBody(expense *models.Expense) Body {
	return Body{
		ID: expense.Id().String(),
		Amount: Money{
			Amount:   expense.Amount().Amount(),
			Currency: expense.Amount().Currency(),
		},
		ExpenseDate: expense.ExpenseDate().Format(DateFormat),
		Description: expense.Description(),
		ExpenseType: TypeBody{
			ID:   expense.ExpenseType().Id().String(),
			Name: expense.ExpenseType().Name(),
		},
	}
}

type AddExpenseRequest struct {
	Amount      Money                             `json:"amount,omitempty"`
	ExpenseDate string                            `json:"expense_date,omitempty" validate:"required,datetime=2006-01-02"`
	Description string                            `json:"description,omitempty"`
	ExpenseType *AddExpenseRequestExpenseTypeBody `json:"expense_type,omitempty" validate:"required"`
}

type AddExpenseRequestExpenseTypeBody struct {
	ID string `json:"id" validate:"required,uuid"`
}

type SearchInPeriodQueryParams struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02,lteStrDateField=EndDate0x2C2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

type Response struct {
	Expense Body `json:"expense"`
}

type SearchResponse struct {
	Expenses []Body `json:"expenses"`
}

type Body struct {
	ID          string   `json:"id"`
	Amount      Money    `json:"amount"`
	ExpenseDate string   `json:"expense_date"`
	Description string   `json:"description"`
	ExpenseType TypeBody `json:"expense_type"`
}

type TypeBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Money struct {
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"iso4217"`
}
