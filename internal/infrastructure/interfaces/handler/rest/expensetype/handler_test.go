package expensetype_test

import (
	"encoding/json"
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expense"
	expenseTypeService "finfit-backend/internal/domain/services/expensetype"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest/expensetype"
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
	errorResponse = `{"status_code":%d,"msg":"%s","error_detail":"%v","field_errors":%v,"error_code":%d}
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

	requestBody := suite.getAddExpenseRequestBodyFromExpenseType(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := suite.getAddExpenseTypeResponseFromExpenseType(expectedAddedExpenseType)

	addCommand, _ := expenseTypeService.NewAddCommand(expectedAddedExpenseType.Name())
	suite.expenseTypeServiceMock.MockAdd([]interface{}{addCommand},
		[]interface{}{expectedAddedExpenseType, nil}, 1)

	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithEmptyName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("")

	requestBody := suite.getAddExpenseRequestBodyFromExpenseType(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expensetype.FieldValidationErrorMessage, expensetype.FieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]", rest.FieldValidationErrorCode)

	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithBlankName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("   ")

	requestBody := suite.getAddExpenseRequestBodyFromExpenseType(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expensetype.FieldValidationErrorMessage, expensetype.FieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name is a required field\"}]", rest.FieldValidationErrorCode)

	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithTooSmallName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("PR")

	requestBody := suite.getAddExpenseRequestBodyFromExpenseType(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expensetype.FieldValidationErrorMessage, expensetype.FieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name must be at least 3 characters in length\"}]", rest.FieldValidationErrorCode)

	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenAnExpenseTypeToAddWithTooLongName_WhenAdd_ThenReturnStatusBadRequest() {
	expectedAddedExpenseType := models.NewExpenseType("THIS IS A VERY LONG NAME FOR EXPENSE TYPE")

	requestBody := suite.getAddExpenseRequestBodyFromExpenseType(expectedAddedExpenseType)
	c, rec := suite.mockAddExpenseTypeRequest(requestBody)

	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusBadRequest, expensetype.FieldValidationErrorMessage, expensetype.FieldValidationErrorMessage, "[{\"field\":\"Name\",\"message\":\"Name must be a maximum of 32 characters in length\"}]", rest.FieldValidationErrorCode)

	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.Add(c)) {
		assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGetAllSuccess() {
	expectedExpenseTypes := []*models.ExpenseType{
		models.NewExpenseType("test1"),
		models.NewExpenseType("test2"),
	}

	expectedResponseBody := suite.getGetAllExpenseTypeResponseFromExpenseTypes(expectedExpenseTypes)
	suite.expenseTypeServiceMock.MockGetAll([]interface{}{}, []interface{}{expectedExpenseTypes, nil}, 1)
	c, rec := suite.mockGetAllExpenseTypeRequest()
	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.GetAll(c)) {
		suite.expenseTypeServiceMock.AssertCalled(suite.T(), "GetAll")
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), expectedResponseBody, rec.Body.String())
	}
}

func (suite *HandlerTestSuite) TestGivenThatServiceReturnUnexpectedError_whenGetAll_thenReturnErrorResponseWithInternalServerErrorStatus() {
	expectedServiceError := expense.UnexpectedError{Msg: "fail"}
	suite.expenseTypeServiceMock.MockGetAll([]interface{}{}, []interface{}{nil, expectedServiceError}, 1)
	expectedResponseBody := fmt.Sprintf(errorResponse, http.StatusInternalServerError, expensetype.UnexpectedErrorMessage, expectedServiceError.Error(), "[]", 0)

	c, rec := suite.mockGetAllExpenseTypeRequest()
	handler := expensetype.NewHandler(suite.expenseTypeServiceMock, suite.getValidator())

	if assert.NoError(suite.T(), handler.GetAll(c)) {
		suite.expenseTypeServiceMock.AssertCalled(suite.T(), "GetAll")
		assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
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

func (suite *HandlerTestSuite) mockGetAllExpenseTypeRequest() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expense-type", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func (suite *HandlerTestSuite) getAddExpenseRequestBodyFromExpenseType(expenseType *models.ExpenseType) string {
	addExpenseTypeBody := expensetype.AddExpenseTypeRequest{
		Name: expenseType.Name(),
	}

	bodyBytes, _ := json.Marshal(addExpenseTypeBody)
	return string(bodyBytes)
}

func (suite *HandlerTestSuite) getAddExpenseTypeResponseFromExpenseType(expenseType *models.ExpenseType) string {
	response := expensetype.AddExpenseTypeResponse{
		ExpenseType: suite.mapExpenseTypeToExpenseTypeBody(expenseType),
	}

	bodyBytes, _ := json.Marshal(response)
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) getGetAllExpenseTypeResponseFromExpenseTypes(expenseTypes []*models.ExpenseType) string {
	bodyBytes, _ := json.Marshal(suite.mapExpenseTypesToGetAllResponse(expenseTypes))
	return string(bodyBytes) + "\n"
}

func (suite *HandlerTestSuite) mapExpenseTypesToGetAllResponse(expenseTypes []*models.ExpenseType) expensetype.GetAllResponse {
	expenseTypeBodies := []expensetype.Body{}
	for _, expenseType := range expenseTypes {
		expenseTypeBodies = append(expenseTypeBodies, suite.mapExpenseTypeToExpenseTypeBody(expenseType))
	}

	return expensetype.GetAllResponse{ExpenseTypes: expenseTypeBodies}
}

func (suite *HandlerTestSuite) mapExpenseTypeToExpenseTypeBody(expenseType *models.ExpenseType) expensetype.Body {
	return expensetype.Body{
		ID:   expenseType.Id().String(),
		Name: expenseType.Name(),
	}
}

func (suite *HandlerTestSuite) getValidator() fieldvalidation.FieldsValidator {
	validator, _ := fieldvalidation.RegisterFieldsValidator(nil, nil)
	return validator
}
