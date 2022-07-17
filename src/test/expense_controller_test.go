package test

import (
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/use_cases/custom_errors"
	"finfit-backend/src/domain/use_cases/service"
	"finfit-backend/src/infrastructure/interfaces/controller"
	"finfit-backend/src/infrastructure/interfaces/controller/validators"
	"finfit-backend/src/test/mock/repository_mock"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGivenAnExpenseToCreate_WhenCreate_ThenReturnStatusOkWithCreatedExpense(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expectedExpenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expectedCreatedExpense := entities.NewExpenseWithId(1, 100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		expectedExpenseType)
	expectedResponseBody := "{\"id\":1,\"amount\":100.2,\"expense_date\":\"2022-03-15T10:04:05Z\",\"description\":\"Lomitos\",\"expense_type\":{\"id\":1,\"name\":\"Delivery\"}}\n"

	expenseServiceMock.On("Create", service.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(&expectedCreatedExpense, nil)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	}
}

func TestGivenAnExpenseToCreateWithoutDescription_WhenCreate_ThenReturnStatusOkWithCreatedExpense(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expectedExpenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expectedCreatedExpense := entities.NewExpenseWithId(1, 100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"",
		expectedExpenseType)
	expectedResponseBody := "{\"id\":1,\"amount\":100.2,\"expense_date\":\"2022-03-15T10:04:05Z\",\"description\":\"\",\"expense_type\":{\"id\":1,\"name\":\"Delivery\"}}\n"

	expenseServiceMock.On("Create", service.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(&expectedCreatedExpense, nil)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	}
}

func TestGivenAnInvalidExpenseType_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expenseServiceMock.On("Create", service.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(nil, custom_errors.InvalidExpenseTypeError{Msg: "the expense type doesn't exists"})

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	expectedResponseBody := "{\"status_code\":400,\"msg\":\"the expense type doesn't exists\",\"error_detail\":\"the expense type doesn't exists\"}\n"
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnUnexpectedError_WhenCreate_ThenReturnErrorWithInternalServerErrorStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expenseServiceMock.On("Create", service.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(nil, custom_errors.UnexpectedError{Msg: "cagamo fuego"})

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	expectedResponseBody := "{\"status_code\":500,\"msg\":\"unexpected error\",\"error_detail\":\"cagamo fuego\"}\n"
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutAmount_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"CreateExpenseRequest.Amount\",\"message\":\"Amount is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithAmountLowerThanZero_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":-1,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"CreateExpenseRequest.Amount\",\"message\":\"Amount must be greater than 0\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseDate_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_type":{"id":1,"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"CreateExpenseRequest.ExpenseDate\",\"message\":\"ExpenseDate is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseType_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":10.3,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z"}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"CreateExpenseRequest.ExpenseType\",\"message\":\"ExpenseType is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseTypeID_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"CreateExpenseRequest.ExpenseType.ID\",\"message\":\"ID is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseTypeName_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"CreateExpenseRequest.ExpenseType.Name\",\"message\":\"Name is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := controller.NewExpenseController(expenseServiceMock, getValidator())

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func mockCreateExpenseRequest(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func getValidator() validators.FieldsValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, translator)

	return validators.NewGenericFieldsValidator(validate, translator)
}
