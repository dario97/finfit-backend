package test

import (
	"finfit/finfit-backend/domain/entities"
	"finfit/finfit-backend/domain/use_cases/custom_errors"
	service2 "finfit/finfit-backend/domain/use_cases/service"
	"finfit/finfit-backend/interfaces/controller"
	"finfit/finfit-backend/test/mock/repository_mock"
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

	expenseServiceMock.On("Create", service2.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(&expectedCreatedExpense, nil)

	handler := controller.NewExpenseController(expenseServiceMock)

	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	}
}

func TestGivenThatExpenseServiceReturnInvalidExpenseTypeError_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expenseServiceMock.On("Create", service2.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(nil, custom_errors.InvalidExpenseTypeError{Msg: "the expense type doesn't exists"})

	handler := controller.NewExpenseController(expenseServiceMock)

	handler.Create(c)

	expectedResponseBody := "{\"status_code\":400,\"msg\":\"the expense type doesn't exists\",\"error_detail\":\"the expense type doesn't exists\"}\n"
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenThatExpenseServiceReturnInternalError_WhenCreate_ThenReturnErrorWithInternalServerErrorStatus(t *testing.T) {
	expenseServiceMock := repository_mock.NewExpenseServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expenseServiceMock.On("Create", service2.NewCreateExpenseCommand(100.2,
		time.Date(2022, time.March, 15, 10, 4, 5, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"),
	)).Return(nil, custom_errors.InternalError{Msg: "cagamo fuego"})

	handler := controller.NewExpenseController(expenseServiceMock)

	handler.Create(c)

	expectedResponseBody := "{\"status_code\":500,\"msg\":\"unexpected error\",\"error_detail\":\"cagamo fuego\"}\n"
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

//
//func TestGivenAnExpenseBodyWithoutAmount_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
//	expenseJson := `{"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z","expense_type":{"id":1,"name":"Delivery"}}`
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(expenseJson))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	service := service2.NewExpenseService(repository_mock.NewExpenseRepositoryMock(errors.New("db error")))
//	c := e.NewContext(req, rec)
//	handler := controller.NewExpenseController(service)
//
//	handler.Create(c)
//
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//}
//
//func TestGivenAnExpenseBodyWithoutExpenseType_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
//	expenseJson := `{"amount":10.3,"description":"Lomitos","expense_date":"2022-03-15T10:04:05Z"}`
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(expenseJson))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	service := service2.NewExpenseService(repository_mock.NewExpenseRepositoryMock(errors.New("db error")))
//	c := e.NewContext(req, rec)
//	handler := controller.NewExpenseController(service)
//
//	handler.Create(c)
//
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//}

func mockCreateExpenseRequest(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
