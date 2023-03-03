package expensetype

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"finfit-backend/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	repositoryMock *RepositoryMock
	service        Service
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.repositoryMock = NewRepositoryMock()
	suite.service = NewService(suite.repositoryMock)
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
	expectedExpenseType := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByID([]interface{}{expectedExpenseType.Id}, []interface{}{expectedExpenseType, nil}, 1)

	actualExpenseType, err := suite.service.GetById(expectedExpenseType.Id)

	require.NoError(suite.T(), err)
	suite.assertEqualsExpenseType(expectedExpenseType, actualExpenseType)
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFails_whenGetById_thenReturnError() {
	id := uuid.New()
	suite.repositoryMock.MockGetByID([]interface{}{id}, []interface{}{nil, errors.New("fail")}, 1)

	actualExpenseType, err := suite.service.GetById(id)

	require.ErrorAs(suite.T(), err, &UnexpectedError{})
	require.Nil(suite.T(), actualExpenseType)
}

func (suite *ServiceTestSuite) TestGivenAnExpenseTypeToAddAndExpenseTypeNotExists_whenAdd_thenReturnAddedExpenseType() {
	expectedExpenseType := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name}, []interface{}{nil, nil}, 1)
	suite.repositoryMock.MockAdd([]interface{}{expectedExpenseType}, []interface{}{expectedExpenseType, nil}, 1)

	addedExpenseType, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.NoError(suite.T(), err)
	suite.assertEqualsExpenseType(expectedExpenseType, addedExpenseType)
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGivenThatExpenseTypeAlreadyExists_whenAdd_thenReturnAddedExpenseType() {
	expectedExpenseType := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name}, []interface{}{expectedExpenseType, nil}, 1)

	addedExpenseType, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.NoError(suite.T(), err)
	suite.assertEqualsExpenseType(expectedExpenseType, addedExpenseType)
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFailsGivingExpenseTypeByName_whenAdd_thenReturnError() {
	expectedExpenseType := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name}, []interface{}{nil, errors.New("fail")}, 1)

	_, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.ErrorAs(suite.T(), err, &UnexpectedError{})
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGivenThatRepositoryFailsAddingExpenseType_whenAdd_thenReturnError() {
	expectedExpenseType := models.NewExpenseType("Servicios")
	suite.repositoryMock.MockGetByName([]interface{}{expectedExpenseType.Name}, []interface{}{nil, nil}, 1)
	suite.repositoryMock.MockAdd([]interface{}{expectedExpenseType}, []interface{}{nil, errors.New("fail")}, 1)

	_, err := suite.service.Add(suite.buildAddCommandFromExpenseType(expectedExpenseType))

	require.ErrorAs(suite.T(), err, &UnexpectedError{})
	suite.repositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) assertEqualsExpenseType(expected *models.ExpenseType, actual *models.ExpenseType) {
	require.Equal(suite.T(), expected.Id, actual.Id)
	require.Equal(suite.T(), expected.Name, actual.Name)
}

func (suite *ServiceTestSuite) buildAddCommandFromExpenseType(expenseType *models.ExpenseType) *AddCommand {
	return &AddCommand{name: expenseType.Name}
}
