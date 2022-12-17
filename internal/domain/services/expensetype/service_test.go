package expensetype

import (
	"errors"
	"finfit-backend/internal/domain/models"
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
}

func (suite *ServiceTestSuite) TearDownTest() {
	suite.repositoryMock.ExpectedCalls = nil
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

func (suite *ServiceTestSuite) assertEqualsExpenseType(expected *models.ExpenseType, actual *models.ExpenseType) {
	require.Equal(suite.T(), expected.Id, actual.Id)
	require.Equal(suite.T(), expected.Name, actual.Name)
}
