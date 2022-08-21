package add

import (
	"finfit-backend/src/domain/use_cases/custom_errors"
	validator2 "finfit-backend/src/pkg/validator"
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
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expectedCreatedExpense := createdExpense{
		id:          1,
		amount:      100.2,
		expenseDate: time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		description: "Lomitos",
		expenseType: expenseType{
			id:   1,
			name: "Delivery",
		},
	}
	expectedResponseBody := "{\"id\":1,\"amount\":100.2,\"expense_date\":\"2022-03-15\",\"description\":\"Lomitos\",\"expense_type\":{\"id\":1,\"name\":\"Delivery\"}}\n"

	expenseServiceMock.On("Add", command{
		amount:      expectedCreatedExpense.amount,
		expenseDate: expectedCreatedExpense.expenseDate,
		description: expectedCreatedExpense.description,
		expenseType: expectedCreatedExpense.expenseType,
	}).Return(&expectedCreatedExpense, nil)

	handler := NewHandler(expenseServiceMock, getValidator())

	if assert.NoError(t, handler.Add(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	}
}

func TestGivenAnExpenseToCreateWithoutDescription_WhenCreate_ThenReturnStatusOkWithCreatedExpense(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"expense_date":"2022-03-15","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expectedCreatedExpense := createdExpense{
		id:          1,
		amount:      100.2,
		expenseDate: time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		description: "",
		expenseType: expenseType{
			id:   1,
			name: "Delivery",
		},
	}
	expectedResponseBody := "{\"id\":1,\"amount\":100.2,\"expense_date\":\"2022-03-15\",\"description\":\"\",\"expense_type\":{\"id\":1,\"name\":\"Delivery\"}}\n"

	expenseServiceMock.On("Add", command{
		amount:      expectedCreatedExpense.amount,
		expenseDate: expectedCreatedExpense.expenseDate,
		description: "",
		expenseType: expectedCreatedExpense.expenseType,
	}).Return(&expectedCreatedExpense, nil)

	handler := NewHandler(expenseServiceMock, getValidator())

	if assert.NoError(t, handler.Add(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	}
}

func TestGivenAnInvalidExpenseType_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expenseServiceMock.On("Add", command{
		amount:      100.2,
		expenseDate: time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		description: "Lomitos",
		expenseType: expenseType{
			id:   1,
			name: "Delivery",
		},
	}).Return(nil, custom_errors.InvalidExpenseTypeError{Msg: "the expenseDBModel type doesn't exists"})

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	expectedResponseBody := "{\"status_code\":400,\"msg\":\"the expenseDBModel type doesn't exists\",\"error_detail\":\"the expenseDBModel type doesn't exists\"}\n"
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnUnexpectedError_WhenCreate_ThenReturnErrorWithInternalServerErrorStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":1,"name":"Delivery"}}`
	c, rec := mockCreateExpenseRequest(requestBody)

	expenseServiceMock.On("Add", command{
		amount:      100.2,
		expenseDate: time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		description: "Lomitos",
		expenseType: expenseType{
			id:   1,
			name: "Delivery",
		},
	}).Return(nil, custom_errors.UnexpectedError{Msg: "cagamo fuego"})

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	expectedResponseBody := "{\"status_code\":500,\"msg\":\"unexpected error\",\"error_detail\":\"cagamo fuego\"}\n"
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutAmount_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":1,"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"Amount\",\"message\":\"Amount is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithAmountLowerThanZero_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":-1,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":1,"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"Amount\",\"message\":\"Amount must be greater than 0\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseDate_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_type":{"id":1,"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"ExpenseDate\",\"message\":\"ExpenseDate is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseType_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":10.3,"description":"Lomitos","expense_date":"2022-03-15"}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"ExpenseType\",\"message\":\"ExpenseType is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseTypeID_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"name":"Delivery"}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"ID\",\"message\":\"ID is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func TestGivenAnExpenseWithoutExpenseTypeName_WhenCreate_ThenReturnErrorWithBadRequestStatus(t *testing.T) {
	expenseServiceMock := NewServiceMock()
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":1}}`
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]}\n"
	c, rec := mockCreateExpenseRequest(requestBody)

	handler := NewHandler(expenseServiceMock, getValidator())

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}

func mockCreateExpenseRequest(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenseDBModel", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func getValidator() validator2.FieldsValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, translator)

	return validator2.NewGenericFieldsValidator(validate, translator)
}
