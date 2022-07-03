package controller

import (
	"errors"
	"finfit/finfit-backend/src/domain/entities"
	"finfit/finfit-backend/src/domain/use_cases/custom_errors"
	"finfit/finfit-backend/src/domain/use_cases/service"
	"finfit/finfit-backend/src/interfaces/controller/dto"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"net/http"
)

const (
	fieldValidationErrorMessage = "some fields are invalid"
	bodyWasInvalidErrorMessage  = "body was invalid"
	unexpectedErrorMessage      = "unexpected error"
)

type ExpenseController interface {
	GetById(context echo.Context) error
	Search(context echo.Context) error
	Create(context echo.Context) error
	DeleteById(context echo.Context) error
	Update(context echo.Context) error
}

type expenseController struct {
	expenseService service.ExpenseService
	validator      *validator.Validate
}

func NewExpenseController(service service.ExpenseService) ExpenseController {
	return expenseController{expenseService: service, validator: validator.New()}
}

func (e expenseController) Search(context echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (e expenseController) GetById(context echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (e expenseController) Create(context echo.Context) error {
	createExpenseRequest := new(dto.CreateExpenseRequest)

	if err := context.Bind(createExpenseRequest); err != nil {
		return e.buildErrorResponse(context, http.StatusBadRequest, bodyWasInvalidErrorMessage, err.Error())
	}

	if err := e.validator.Struct(createExpenseRequest); err != nil {
		fieldValidationErrors := e.buildFieldValidationErrors(err.(validator.ValidationErrors))
		return e.buildErrorResponse(context, http.StatusBadRequest, fieldValidationErrorMessage, fieldValidationErrors)
	}

	expenseType := entities.NewExpenseTypeWithId(createExpenseRequest.ExpenseType.ID, createExpenseRequest.ExpenseType.Name)

	createdExpense, err := e.expenseService.Create(service.NewCreateExpenseCommand(createExpenseRequest.Amount,
		createExpenseRequest.ExpenseDate,
		createExpenseRequest.Description,
		expenseType))

	if err != nil {
		return e.manageServiceError(context, err)
	}

	return context.JSON(http.StatusCreated, dto.NewExpenseResponseFromExpense(*createdExpense))
}

func (e expenseController) DeleteById(context echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (e expenseController) Update(context echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (e expenseController) buildErrorResponse(ctx echo.Context, statusCode int, errorMessage string, errorDetail interface{}) error {
	errorResponse := errorResponse{StatusCode: statusCode, Msg: errorMessage, ErrorDetail: errorDetail}
	return ctx.JSON(statusCode, errorResponse)
}

func (e expenseController) buildFieldValidationErrors(fieldErrors []validator.FieldError) []FieldValidationError {
	var fieldValidationErrors []FieldValidationError
	for _, validationError := range fieldErrors {
		fieldValidationError := FieldValidationError{Field: validationError.Namespace(), ValidationResult: validationError.Tag()}
		fieldValidationErrors = append(fieldValidationErrors, fieldValidationError)
	}

	return fieldValidationErrors
}

func (e expenseController) manageServiceError(ctx echo.Context, err error) error {
	if errors.As(err, &custom_errors.InvalidExpenseTypeError{}) {
		return e.buildErrorResponse(ctx, http.StatusBadRequest, err.Error(), err.Error())
	} else {
		return e.buildErrorResponse(ctx, http.StatusInternalServerError, unexpectedErrorMessage, err.Error())
	}
}
