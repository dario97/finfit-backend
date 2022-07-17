package controller

import (
	"errors"
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/use_cases/custom_errors"
	"finfit-backend/src/domain/use_cases/service"
	dto2 "finfit-backend/src/infrastructure/interfaces/controller/dto"
	"finfit-backend/src/infrastructure/interfaces/controller/validators"
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
	expenseService  service.ExpenseService
	fieldsValidator validators.FieldsValidator
}

func NewExpenseController(service service.ExpenseService, fieldsValidator validators.FieldsValidator) ExpenseController {
	return expenseController{expenseService: service, fieldsValidator: fieldsValidator}
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
	createExpenseRequest := new(dto2.CreateExpenseRequest)

	if err := context.Bind(createExpenseRequest); err != nil {
		return e.buildErrorResponse(context, http.StatusBadRequest, bodyWasInvalidErrorMessage, err.Error())
	}

	if fieldValidationErrors := e.fieldsValidator.ValidateFields(createExpenseRequest); len(fieldValidationErrors) > 0 {
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

	return context.JSON(http.StatusCreated, dto2.NewExpenseResponseFromExpense(*createdExpense))
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

func (e expenseController) manageServiceError(ctx echo.Context, err error) error {
	if errors.As(err, &custom_errors.InvalidExpenseTypeError{}) {
		return e.buildErrorResponse(ctx, http.StatusBadRequest, err.Error(), err.Error())
	} else {
		return e.buildErrorResponse(ctx, http.StatusInternalServerError, unexpectedErrorMessage, err.Error())
	}
}
