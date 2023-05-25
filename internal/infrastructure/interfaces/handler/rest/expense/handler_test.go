package expense_test

import (
	"encoding/json"
	"finfit-backend/internal/domain/models"
	expenseService "finfit-backend/internal/domain/services/expense"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest/expense"
	"finfit-backend/pkg"
	"finfit-backend/pkg/fieldvalidation"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	errorResponse = `{"status_code":%d,"msg":"%s","error_detail":"%v","field_errors":%v,"error_code":%d}
`
)

type HandlerTestSuite struct {
	suite.Suite
	expenseServiceMock *expenseService.ServiceMock
}

func (suite *HandlerTestSuite) SetupSuite() {
	suite.expenseServiceMock = expenseService.NewServiceMock()
	suite.patchNewUUIDMethod()
}

func (suite *HandlerTestSuite) patchNewUUIDMethod() {
	id := uuid.New()
	pkg.NewUUID = func() uuid.UUID {
		return id
	}
}

func (suite *HandlerTestSuite) TearDownSuite() {
	suite.unpatchNewUUIDMethod()
}

func (suite *HandlerTestSuite) unpatchNewUUIDMethod() {
	pkg.NewUUID = uuid.New
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestGivenAnExpenseToCreate_WhenAdd_ThenReturnStatusOkWithCreatedExpense() {
	expectedCreatedExpense := models.NewExpense(
		models.NewMoney(100.2, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedCreatedExpense)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := suite.getAddExpenseResponseFromExpense(expectedCreatedExpense)

	addCommand, _ := expenseService.NewAddCommand(expectedCreatedExpense.Amount().Amount(),
		expectedCreatedExpense.Amount().Currency(),
		expectedCreatedExpense.ExpenseDate(),
		expectedCreatedExpense.Description(),
		expectedCreatedExpense.ExpenseType().Id())
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{expectedCreatedExpense, nil}, 1)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseToCreateWithoutDescription_WhenAdd_ThenReturnStatusOkWithCreatedExpense() {
	expectedCreatedExpense := models.NewExpense(models.NewMoney(100.2, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedCreatedExpense)
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := suite.getAddExpenseResponseFromExpense(expectedCreatedExpense)

	addCommand, _ := expenseService.NewAddCommand(expectedCreatedExpense.Amount().Amount(),
		expectedCreatedExpense.Amount().Currency(),
		expectedCreatedExpense.ExpenseDate(),
		expectedCreatedExpense.Description(),
		expectedCreatedExpense.ExpenseType().Id())
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{expectedCreatedExpense, nil}, 1)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnInvalidExpenseType_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(100.2, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	addCommand, _ := expenseService.NewAddCommand(expenseToCreate.Amount().Amount(),
		expenseToCreate.Amount().Currency(),
		expenseToCreate.ExpenseDate(),
		expenseToCreate.Description(),
		expenseToCreate.ExpenseType().Id())

	serviceErr := expenseService.InvalidExpenseTypeError{Msg: "the expense type doesn't exists"}
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{nil, serviceErr}, 1)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, serviceErr.Error(), serviceErr.Error(), "[]", 0)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnUnexpectedError_WhenAdd_ThenReturnErrorWithInternalServerErrorStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(100.2, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	addCommand, _ := expenseService.NewAddCommand(expenseToCreate.Amount().Amount(),
		expenseToCreate.Amount().Currency(),
		expenseToCreate.ExpenseDate(),
		expenseToCreate.Description(),
		expenseToCreate.ExpenseType().Id())
	serviceErr := expenseService.UnexpectedError{Msg: "cagamo fuego"}
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{nil, serviceErr}, 1)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusInternalServerError, expense.UnexpectedErrorMessage, serviceErr.Error(), "[]", 0)
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutAmount_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(0, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expense.FieldValidationErrorMessage, expense.FieldValidationErrorMessage, `[{"field":"Amount","message":"Amount is a required field"}]`, rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithAmountLowerThanZero_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(-1, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		expense.FieldValidationErrorMessage,
		expense.FieldValidationErrorMessage,
		`[{"field":"Amount","message":"Amount must be greater than 0"}]`,
		rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseDate_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(100, "ARS"),
		time.Time{},
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := fmt.Sprintf(`{"amount":{"amount":%f,"currency":"%s"},"description":"%s","expense_type":{"id":"%s"}}`,
		expenseToCreate.Amount().Amount(),
		expenseToCreate.Amount().Currency(),
		expenseToCreate.Description(),
		expenseToCreate.ExpenseType().Id().String())

	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		expense.FieldValidationErrorMessage,
		expense.FieldValidationErrorMessage,
		`[{"field":"ExpenseDate","message":"ExpenseDate is a required field"}]`,
		rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithBadFormattedExpenseDate_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(100, "ARS"),
		time.Time{}, "Lomitos", models.NewExpenseType("Delivery"))

	requestBody := fmt.Sprintf(`{"amount":{"amount":%f,"currency":"%s"},"expense_date":"%s","description":"%s","expense_type":{"id":"%s"}}`,
		expenseToCreate.Amount().Amount(),
		expenseToCreate.Amount().Currency(),
		"12-2013-12",
		expenseToCreate.Description(),
		expenseToCreate.ExpenseType().Id().String())

	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		expense.FieldValidationErrorMessage,
		expense.FieldValidationErrorMessage,
		`[{"field":"ExpenseDate","message":"ExpenseDate does not match the 2006-01-02 format"}]`,
		rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseType_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(models.NewMoney(10.2, "ARS"),
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", nil)

	requestBody := fmt.Sprintf(`{"amount":{"amount":%f,"currency":"%s"},"expense_date":"%s","description":"%s"}`,
		expenseToCreate.Amount().Amount(),
		expenseToCreate.Amount().Currency(),
		"2013-02-01",
		expenseToCreate.Description())
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		expense.FieldValidationErrorMessage,
		expense.FieldValidationErrorMessage,
		`[{"field":"ExpenseType","message":"ExpenseType is a required field"}]`,
		rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseTypeID_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	requestBody := `{"amount":{"amount":100.2,"currency":"ARS"},"description":"Lomitos","expense_date":"2022-03-15","expense_type":{}}`
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":\"some fields are invalid\",\"field_errors\":[{\"field\":\"ID\",\"message\":\"ID is a required field\"}],\"error_code\":1}\n"

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithNoUIIDExpenseTypeID_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	requestBody := `{"amount":{"amount":100.2,"currency":"ARS"},"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":"fruta-uuid"}}`
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":\"some fields are invalid\",\"field_errors\":[{\"field\":\"ID\",\"message\":\"ID must be a valid UUID\"}],\"error_code\":1}\n"

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAPeriod_WhenSearchInPeriod_ThenReturnStatusOkWithListOfExpenses() {
	expectedExpensesToReturn := suite.getExpenses()
	startDate := time.Date(2022, 5, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)
	searchInPeriodCommand, _ := expenseService.NewSearchInPeriodCommand(startDate, endDate)
	suite.expenseServiceMock.MockSearchInPeriod([]interface{}{searchInPeriodCommand}, []interface{}{expectedExpensesToReturn, nil}, 1)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(expense.DateFormat), endDate.Format(expense.DateFormat)))

	expectedResponseBody := suite.getSearchResponseBodyFromExpenses(expectedExpensesToReturn)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.SearchInPeriod(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenThatServiceFails_WhenSearchInPeriod_ThenReturnStatusInternalServerError() {
	startDate := time.Date(2022, 5, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)
	searchInPeriodCommand, _ := expenseService.NewSearchInPeriodCommand(startDate, endDate)
	expectedServiceError := expenseService.UnexpectedError{Msg: "fail getting expenses"}
	suite.expenseServiceMock.MockSearchInPeriod([]interface{}{searchInPeriodCommand}, []interface{}{nil, expectedServiceError}, 1)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(expense.DateFormat), endDate.Format(expense.DateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusInternalServerError, expense.UnexpectedErrorMessage, expectedServiceError.Error(), "[]", 0)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatStartDateParamNotExists_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("end_date=%s", endDate.Format(expense.DateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expense.FieldValidationErrorMessage, expense.FieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate is a required field\"}]", rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatEndDateParamNotExists_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s", startDate.Format(expense.DateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expense.FieldValidationErrorMessage, expense.FieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate must be before or equal to EndDate\"},{\"field\":\"EndDate\",\"message\":\"EndDate is a required field\"}]", rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatStartDateParamHasBadFormat_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format("02-01-2006"), endDate.Format(expense.DateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expense.FieldValidationErrorMessage, expense.FieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate does not match the 2006-01-02 format\"}]", rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatEndDateParamHasBadFormat_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(expense.DateFormat), endDate.Format("02-01-2006")))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expense.FieldValidationErrorMessage, expense.FieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate must be before or equal to EndDate\"},{\"field\":\"EndDate\",\"message\":\"EndDate does not match the 2006-01-02 format\"}]", rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatStartDateIsGreaterThanEndDate_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 9, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(expense.DateFormat), endDate.Format(expense.DateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expense.FieldValidationErrorMessage, expense.FieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate must be before or equal to EndDate\"}]", rest.FieldValidationErrorCode)

	handler := expense.NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) getExpenses() []*models.Expense {
	return []*models.Expense{models.NewExpense(models.NewMoney(100.2, "ARS"),
		time.Date(2022, time.May, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery")),
		models.NewExpense(models.NewMoney(100.2, "ARS"),
			time.Date(2022, time.September, 15, 0, 0, 0, 0, time.UTC),
			"Lomitos", models.NewExpenseType("Delivery"))}
}

func (suite *HandlerTestSuite) getValidator() fieldvalidation.FieldsValidator {
	validator, _ := fieldvalidation.RegisterFieldsValidator(nil, nil)
	return validator
}

func (suite *HandlerTestSuite) getAddExpenseResponseFromExpense(domainExpense *models.Expense) string {
	response := expense.Response{Expense: expense.Body{
		ID: domainExpense.Id().String(),
		Amount: expense.Money{
			Amount:   domainExpense.Amount().Amount(),
			Currency: domainExpense.Amount().Currency(),
		},
		ExpenseDate: domainExpense.ExpenseDate().Format(expense.DateFormat),
		Description: domainExpense.Description(),
		ExpenseType: expense.TypeBody{
			ID:   domainExpense.ExpenseType().Id().String(),
			Name: domainExpense.ExpenseType().Name(),
		},
	}}

	bodyBytes, _ := json.Marshal(response)
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) getAddExpenseRequestBodyFromExpense(domainExpense *models.Expense) string {
	addExpenseBody := expense.AddExpenseRequest{
		Amount: expense.Money{
			Amount:   domainExpense.Amount().Amount(),
			Currency: domainExpense.Amount().Currency(),
		},
		ExpenseDate: domainExpense.ExpenseDate().Format(expense.DateFormat),
		Description: domainExpense.Description(),
		ExpenseType: &expense.AddExpenseRequestExpenseTypeBody{
			ID: domainExpense.ExpenseType().Id().String(),
		},
	}

	bodyBytes, _ := json.Marshal(addExpenseBody)
	return string(bodyBytes)
}

func (suite *HandlerTestSuite) mockAddExpenseRequest(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func (suite *HandlerTestSuite) mockSearchInPeriodRequest(queryParams string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/expense/search?%s",
			queryParams),
		nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func (suite *HandlerTestSuite) getSearchResponseBodyFromExpenses(expenses []*models.Expense) string {
	expenseBodies := []expense.Body{}
	for _, expense := range expenses {
		expenseBodies = append(expenseBodies, suite.mapExpenseToExpenseBody(expense))
	}

	response := expense.SearchResponse{Expenses: expenseBodies}
	bodyBytes, _ := json.Marshal(response)
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) mapExpenseToExpenseBody(domainExpense *models.Expense) expense.Body {
	return expense.Body{
		ID: domainExpense.Id().String(),
		Amount: expense.Money{
			Amount:   domainExpense.Amount().Amount(),
			Currency: domainExpense.Amount().Currency(),
		},
		ExpenseDate: domainExpense.ExpenseDate().Format(expense.DateFormat),
		Description: domainExpense.Description(),
		ExpenseType: expense.TypeBody{
			ID:   domainExpense.ExpenseType().Id().String(),
			Name: domainExpense.ExpenseType().Name(),
		},
	}
}
