package expensetype_test

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expensetype"
	"finfit-backend/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	repositoryMock *expensetype.RepositoryMock
	service        expensetype.Service
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.repositoryMock = expensetype.NewRepositoryMock()
	suite.service = expensetype.NewService(suite.repositoryMock)
	suite.patchUUIDFunction()
}

func (suite *ServiceTestSuite) patchUUIDFunction() {
	id := uuid.New()
	pkg.NewUUID = func() uuid.UUID {
		return id
	}
}

func (suite *ServiceTestSuite) TearDownTest() {
	suite.repositoryMock.ExpectedCalls = nil
}

func (suite *ServiceTestSuite) TearDownSuite() {
	pkg.NewUUID = uuid.New
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestGivenAnID_whenGetById_thenReturnExpenseType() {
	expectedExpenseType, _ := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByID([]interface{}{expectedExpenseType.Id()}, []interface{}{expectedExpenseType, nil}, 1)

	actualExpenseType, err := suite.service.GetById(expectedExpenseType.Id())

	require.NoError(suite.T(), err)
	suite.assertEqualsExpenseType(expectedExpenseType, actualExpenseType)
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFails_whenGetById_thenReturnError() {
	id := uuid.New()
	suite.repositoryMock.MockGetByID([]interface{}{id}, []interface{}{nil, errors.New("fail")}, 1)

	actualExpenseType, err := suite.service.GetById(id)

	require.ErrorAs(suite.T(), err, &expensetype.UnexpectedError{})
	require.Nil(suite.T(), actualExpenseType)
}

func (suite *ServiceTestSuite) TestGivenAnExpenseTypeToAddAndExpenseTypeNotExists_whenAdd_thenReturnAddedExpenseType() {
	expectedExpenseType, _ := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name()}, []interface{}{nil, nil}, 1)
	suite.repositoryMock.MockAdd([]interface{}{expectedExpenseType}, []interface{}{expectedExpenseType, nil}, 1)

	addedExpenseType, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.NoError(suite.T(), err)
	suite.assertEqualsExpenseType(expectedExpenseType, addedExpenseType)
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGivenThatExpenseTypeAlreadyExists_whenAdd_thenReturnAddedExpenseType() {
	expectedExpenseType, _ := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name()}, []interface{}{expectedExpenseType, nil}, 1)

	addedExpenseType, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.NoError(suite.T(), err)
	suite.assertEqualsExpenseType(expectedExpenseType, addedExpenseType)
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFailsGivingExpenseTypeByName_whenAdd_thenReturnError() {
	expectedExpenseType, _ := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name()}, []interface{}{nil, errors.New("fail")}, 1)

	_, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.ErrorAs(suite.T(), err, &expensetype.UnexpectedError{})
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFailsAddingExpenseType_whenAdd_thenReturnError() {
	expectedExpenseType, _ := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name()}, []interface{}{nil, nil}, 1)
	suite.repositoryMock.MockAdd([]interface{}{expectedExpenseType}, []interface{}{nil, errors.New("fail")}, 1)

	_, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.ErrorAs(suite.T(), err, &expensetype.UnexpectedError{})
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAll_Success() {
	expectedExpenseTypes := suite.getExpenseTypes()
	suite.repositoryMock.MockGetAll([]interface{}{}, []interface{}{expectedExpenseTypes, nil}, 1)

	actualExpenseTypes, err := suite.service.GetAll()

	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), actualExpenseTypes, 2)
	assert.Equal(suite.T(), actualExpenseTypes[0].Id(), expectedExpenseTypes[0].Id())
	assert.Equal(suite.T(), actualExpenseTypes[0].Name(), expectedExpenseTypes[0].Name())
	assert.Equal(suite.T(), actualExpenseTypes[1].Id(), expectedExpenseTypes[1].Id())
	assert.Equal(suite.T(), actualExpenseTypes[1].Name(), expectedExpenseTypes[1].Name())
}

func (suite *ServiceTestSuite) getExpenseTypes() []*models.ExpenseType {
	return []*models.ExpenseType{suite.getExpenseType1(), suite.getExpenseType2()}
}

func (suite *ServiceTestSuite) getExpenseType1() *models.ExpenseType {
	expenseType, _ := models.NewExpenseType("Food")
	return expenseType
}

func (suite *ServiceTestSuite) getExpenseType2() *models.ExpenseType {
	expenseType, _ := models.NewExpenseType("Travel")
	return expenseType
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFails_whenGetAll_thenReturnError() {
	suite.repositoryMock.MockGetAll([]interface{}{}, []interface{}{nil, errors.New("fail")}, 1)

	expenseTypes, err := suite.service.GetAll()

	assert.NotNil(suite.T(), err)
	assert.ErrorAs(suite.T(), err, &expensetype.UnexpectedError{})
	assert.Nil(suite.T(), expenseTypes)
}

func (suite *ServiceTestSuite) assertEqualsExpenseType(expected *models.ExpenseType, actual *models.ExpenseType) {
	require.Equal(suite.T(), expected.Id(), actual.Id())
	require.Equal(suite.T(), expected.Name(), actual.Name())
}

func (suite *ServiceTestSuite) buildAddCommandFromExpenseType(expenseType *models.ExpenseType) *expensetype.AddCommand {
	command, _ := expensetype.NewAddCommand(expenseType.Name())
	return command
}
