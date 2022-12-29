package expense

import (
	"encoding/json"
	"finfit-backend/internal/domain/models"
	expenseService "finfit-backend/internal/domain/services/expense"
	"finfit-backend/pkg"
	"finfit-backend/pkg/fieldvalidation"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
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
	errorResponse = `{"status_code":%d,"msg":"%s","error_detail":%v}
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
	expectedCreatedExpense := models.NewExpense(100.2,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedCreatedExpense)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := suite.getAddExpenseResponseFromExpense(expectedCreatedExpense)

	addCommand, _ := expenseService.NewAddCommand(expectedCreatedExpense.Amount,
		expectedCreatedExpense.ExpenseDate,
		expectedCreatedExpense.Description,
		expectedCreatedExpense.ExpenseType)
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{expectedCreatedExpense, nil}, 1)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseToCreateWithoutDescription_WhenAdd_ThenReturnStatusOkWithCreatedExpense() {
	expectedCreatedExpense := models.NewExpense(100.2,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedCreatedExpense)
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := suite.getAddExpenseResponseFromExpense(expectedCreatedExpense)

	addCommand, _ := expenseService.NewAddCommand(expectedCreatedExpense.Amount,
		expectedCreatedExpense.ExpenseDate,
		expectedCreatedExpense.Description,
		expectedCreatedExpense.ExpenseType)
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{expectedCreatedExpense, nil}, 1)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnInvalidExpenseType_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(100.2,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	addCommand, _ := expenseService.NewAddCommand(expenseToCreate.Amount,
		expenseToCreate.ExpenseDate,
		expenseToCreate.Description,
		expenseToCreate.ExpenseType)

	serviceErr := expenseService.InvalidExpenseTypeError{Msg: "the expense type doesn't exists"}
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{nil, serviceErr}, 1)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, serviceErr.Error(), fmt.Sprintf(`"%s"`, serviceErr.Error()))
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnUnexpectedError_WhenAdd_ThenReturnErrorWithInternalServerErrorStatus() {
	expenseToCreate := models.NewExpense(100.2,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	addCommand, _ := expenseService.NewAddCommand(expenseToCreate.Amount,
		expenseToCreate.ExpenseDate,
		expenseToCreate.Description,
		expenseToCreate.ExpenseType)
	serviceErr := expenseService.UnexpectedError{Msg: "cagamo fuego"}
	suite.expenseServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{nil, serviceErr}, 1)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusInternalServerError, unexpectedErrorMessage, fmt.Sprintf(`"%s"`, serviceErr.Error()))
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutAmount_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(0,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, `[{"field":"Amount","message":"Amount is a required field"}]`)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithAmountLowerThanZero_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(-1,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expenseToCreate)
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		fieldValidationErrorMessage,
		`[{"field":"Amount","message":"Amount must be greater than 0"}]`)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseDate_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(100,
		time.Time{},
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := fmt.Sprintf(`{"amount":%f,"description":"%s","expense_type":{"id":"%s","name":"%s"}}`,
		expenseToCreate.Amount,
		expenseToCreate.Description,
		expenseToCreate.ExpenseType.Id.String(),
		expenseToCreate.ExpenseType.Name)

	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		fieldValidationErrorMessage,
		`[{"field":"ExpenseDate","message":"ExpenseDate is a required field"}]`)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithBadFormattedExpenseDate_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(100,
		time.Time{},
		"Lomitos", models.NewExpenseType("Delivery"))

	requestBody := fmt.Sprintf(`{"amount":%f,"expense_date":"%s","description":"%s","expense_type":{"id":"%s","name":"%s"}}`,
		expenseToCreate.Amount,
		"12-2013-12",
		expenseToCreate.Description,
		expenseToCreate.ExpenseType.Id.String(),
		expenseToCreate.ExpenseType.Name)

	c, rec := suite.mockAddExpenseRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		fieldValidationErrorMessage,
		`[{"field":"ExpenseDate","message":"ExpenseDate does not match the 2006-01-02 format"}]`)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseType_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	expenseToCreate := models.NewExpense(10.2,
		time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", nil)

	requestBody := fmt.Sprintf(`{"amount":%f,"expense_date":"%s","description":"%s"}`,
		expenseToCreate.Amount,
		"2013-02-01",
		expenseToCreate.Description)
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := fmt.Sprintf(errorResponse,
		http.StatusBadRequest,
		fieldValidationErrorMessage,
		`[{"field":"ExpenseType","message":"ExpenseType is a required field"}]`)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseTypeID_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"name":"Delivery"}}`
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"ID\",\"message\":\"ID is a required field\"}]}\n"

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithNoUIIDExpenseTypeID_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	requestBody := `{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":"fruta-uuid","name":"Delivery"}}`
	c, rec := suite.mockAddExpenseRequest(requestBody)
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"ID\",\"message\":\"ID must be a valid UUID\"}]}\n"

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithoutExpenseTypeName_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	id := pkg.NewUUID()
	requestBody := fmt.Sprintf(`{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":"%s"}}`, id.String())
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]}\n"
	c, rec := suite.mockAddExpenseRequest(requestBody)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.Add(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenAnExpenseWithEmptyExpenseTypeName_WhenAdd_ThenReturnErrorWithBadRequestStatus() {
	id := pkg.NewUUID()
	requestBody := fmt.Sprintf(`{"amount":100.2,"description":"Lomitos","expense_date":"2022-03-15","expense_type":{"id":"%s","name":""}}`, id.String())
	expectedResponseBody := "{\"status_code\":400,\"msg\":\"some fields are invalid\",\"error_detail\":[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]}\n"
	c, rec := suite.mockAddExpenseRequest(requestBody)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

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

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(dateFormat), endDate.Format(dateFormat)))

	expectedResponseBody := suite.getSearchResponseBodyFromExpenses(expectedExpensesToReturn)

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

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

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(dateFormat), endDate.Format(dateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusInternalServerError, unexpectedErrorMessage, fmt.Sprintf(`"%s"`, expectedServiceError.Error()))

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatStartDateParamNotExists_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("end_date=%s", endDate.Format(dateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate is a required field\"}]")

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatEndDateParamNotExists_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s", startDate.Format(dateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"EndDate\",\"message\":\"EndDate is a required field\"}]")

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatStartDateParamHasBadFormat_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format("02-01-2006"), endDate.Format(dateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"StartDate\",\"message\":\"StartDate does not match the 2006-01-02 format\"}]")

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatEndDateParamHasBadFormat_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(dateFormat), endDate.Format("02-01-2006")))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"EndDate\",\"message\":\"EndDate does not match the 2006-01-02 format\"}]")

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) TestGivenThatStartDateIsGreaterThanEndDate_WhenSearchInPeriod_ThenReturnStatusBadRequest() {
	startDate := time.Date(2022, 9, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 8, 13, 0, 0, 0, 0, time.UTC)

	c, rec := suite.mockSearchInPeriodRequest(fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(dateFormat), endDate.Format(dateFormat)))

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"EndDate\",\"message\":\"EndDate does not match the 02-01-2006 format\"}]")

	handler := NewHandler(suite.expenseServiceMock, suite.getValidator())

	handler.SearchInPeriod(c)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
}

func (suite *HandlerTestSuite) getExpenses() []*models.Expense {
	return []*models.Expense{models.NewExpense(100.2,
		time.Date(2022, time.May, 15, 0, 0, 0, 0, time.UTC),
		"Lomitos", models.NewExpenseType("Delivery")),
		models.NewExpense(100.2,
			time.Date(2022, time.September, 15, 0, 0, 0, 0, time.UTC),
			"Lomitos", models.NewExpenseType("Delivery"))}
}

func (suite *HandlerTestSuite) getValidator() fieldvalidation.FieldsValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, translator)

	validate.RegisterValidation(fieldvalidation.LteStrDateFieldValidationTag, fieldvalidation.LteStrDateField)
	return fieldvalidation.NewFieldsValidator(validate, translator)
}

func (suite *HandlerTestSuite) getAddExpenseResponseFromExpense(expense *models.Expense) string {
	response := expenseResponse{Expense: expenseBody{
		ID:          expense.Id.String(),
		Amount:      expense.Amount,
		ExpenseDate: expense.ExpenseDate.Format(dateFormat),
		Description: expense.Description,
		ExpenseType: expenseTypeBody{
			ID:   expense.ExpenseType.Id.String(),
			Name: expense.ExpenseType.Name,
		},
	}}

	bodyBytes, _ := json.Marshal(response)
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) getAddExpenseRequestBodyFromExpense(expense *models.Expense) string {
	addExpenseBody := addExpenseRequest{
		Amount:      expense.Amount,
		ExpenseDate: expense.ExpenseDate.Format(dateFormat),
		Description: expense.Description,
		ExpenseType: &expenseTypeBody{
			ID:   expense.ExpenseType.Id.String(),
			Name: expense.ExpenseType.Name,
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
	expenseBodies := []expenseBody{}
	for _, expense := range expenses {
		expenseBodies = append(expenseBodies, suite.mapExpenseToExpenseBody(expense))
	}

	response := searchResponse{Expenses: expenseBodies}
	bodyBytes, _ := json.Marshal(response)
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) mapExpenseToExpenseBody(expense *models.Expense) expenseBody {
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
