package expensetype

import (
	"encoding/json"
	"finfit-backend/internal/domain/models"
	expenseTypeService "finfit-backend/internal/domain/services/expensetype"
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
)

const (
	errorResponse = `{"status_code":%d,"msg":"%s","error_detail":%v}
`
)

type HandlerTestSuite struct {
	suite.Suite
	expenseTypeServiceMock *expenseTypeService.ServiceMock
}

func (suite *HandlerTestSuite) SetupSuite() {
	suite.expenseTypeServiceMock = expenseTypeService.NewServiceMock()
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

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAdd_WhenAdd_ThenReturnStatusOkWithCreatedExpenseType() {
	expectedAddedExpenseType := models.NewExpenseType("Servicios")

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := suite.getAddExpenseTypeResponseFromExpense(expectedAddedExpenseType)

	addCommand, _ := expenseTypeService.NewAddCommand(expectedAddedExpenseType.Name)
	suite.expenseTypeServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{expectedAddedExpenseType, nil}, 1)

	handler := NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithEmptyName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("")

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]")

	handler := NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithBlankName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("   ")

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]")

	handler := NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithTooSmallName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("PR")

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name must be at least 3 characters in length\"}]")

	handler := NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithTooLongName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("THIS IS A VERY LONG NAME FOR EXPENSE TYPE")

	requestBody := suite.getAddExpenseRequestBodyFromExpense(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, fieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name must be a maximum of 32 characters in length\"}]")

	handler := NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) mockAddExpenseTypeRequest(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expense-type", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func (suite *HandlerTestSuite) getAddExpenseRequestBodyFromExpense(expenseType *models.ExpenseType) string {
	addExpenseTypeBody := addExpenseTypeRequest{
		Name: expenseType.Name,
	}

	bodyBytes, _ := json.Marshal(addExpenseTypeBody)
	return string(bodyBytes)
}

func (suite *HandlerTestSuite) getAddExpenseTypeResponseFromExpense(expenseType *models.ExpenseType) string {
	response := addExpenseTypeResponse{
		ID:   expenseType.Id.String(),
		Name: expenseType.Name,
	}

	bodyBytes, _ := json.Marshal(response)
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) getValidator() fieldvalidation.FieldsValidator {
	validator, _ := fieldvalidation.RegisterFieldsValidator(nil, nil)
	return validator
}
