package expense

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
	"time"
)

type ExpenseServiceTestSuite struct {
	suite.Suite
	expenseRepositoryMock  *RepositoryMock
	expenseTypeServiceMock *expensetype.ServiceMock
	service                Service
}

func (suite *ExpenseServiceTestSuite) SetupSuite() {
	suite.expenseRepositoryMock = NewRepositoryMock()
	suite.expenseTypeServiceMock = expensetype.NewServiceMock()
	suite.service = NewService(suite.expenseRepositoryMock, suite.expenseTypeServiceMock)
	suite.patchUUIDFunction()
}

func (suite *ExpenseServiceTestSuite) patchUUIDFunction() {
	id := uuid.New()
	pkg.NewUUID = func() uuid.UUID {
		return id
	}
}

func (suite *ExpenseServiceTestSuite) TearDownSuite() {
	pkg.NewUUID = uuid.New
}

func (suite *ExpenseServiceTestSuite) TearDownTest() {
	suite.expenseRepositoryMock.ExpectedCalls = nil
	suite.expenseTypeServiceMock.ExpectedCalls = nil
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ExpenseServiceTestSuite))
}

func (suite *ExpenseServiceTestSuite) TestGivenAnExpense_WhenAdd_ThenReturnCreatedExpense() {
	expenseToCreate := getExpense()

	expectedCreatedExpense := &models.Expense{
		Id:          pkg.NewUUID(),
		Amount:      expenseToCreate.Amount,
		ExpenseDate: expenseToCreate.ExpenseDate,
		Description: expenseToCreate.Description,
		ExpenseType: expenseToCreate.ExpenseType,
	}

	suite.expenseRepositoryMock.MockAdd([]interface{}{expenseToCreate}, []interface{}{expectedCreatedExpense, nil}, 1)
	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType.Id}, []interface{}{expenseToCreate.ExpenseType, nil}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), err, "Error must to be nil")
	assertEqualsExpense(suite.T(), expectedCreatedExpense, actualCreatedExpense)
}

func (suite *ExpenseServiceTestSuite) TestGivenThatFailToGetExpenseType_WhenAdd_ThenReturnError() {
	expenseToCreate := getExpense()

	expenseTypeServiceError := errors.New("fail to get expense type")
	expectedError := UnexpectedError{
		Msg: expenseTypeServiceError.Error(),
	}

	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType.Id}, []interface{}{nil, expenseTypeServiceError}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), actualCreatedExpense)
	assert.NotNil(suite.T(), err, "Error must not be nil")
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *ExpenseServiceTestSuite) TestGivenThatExpenseTypeNotExists_WhenAdd_ThenReturnError() {
	expenseToCreate := getExpense()

	expectedError := InvalidExpenseTypeError{
		Msg: "the expense type doesn't exists",
	}

	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType.Id}, []interface{}{nil, nil}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), actualCreatedExpense)
	assert.NotNil(suite.T(), err, "Error must not be nil")
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *ExpenseServiceTestSuite) TestGivenThatSaveExpenseIntoDatabaseFails_WhenAdd_ThenReturnError() {
	expenseToCreate := getExpense()

	repoError := errors.New("fail to save expense")
	expectedError := UnexpectedError{
		Msg: repoError.Error(),
	}

	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType.Id}, []interface{}{expenseToCreate.ExpenseType, nil}, 1)
	suite.expenseRepositoryMock.MockAdd([]interface{}{expenseToCreate}, []interface{}{nil, repoError}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), actualCreatedExpense)
	assert.NotNil(suite.T(), err, "Error must not be nil")
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *ExpenseServiceTestSuite) TestGivenAPeriod_WhenSearchInPeriod_ThenReturnAListOfExpenses() {
	expensesToReturn := suite.getExpenses()

	searchInPeriodCommand, _ := NewSearchInPeriodCommand(
		time.Date(2022, 5, 23, 0, 0, 0, 0, time.Local),
		time.Date(2022, 8, 23, 0, 0, 0, 0, time.Local))

	suite.expenseRepositoryMock.MockSearchInPeriod(
		[]interface{}{searchInPeriodCommand.startDate, searchInPeriodCommand.endDate},
		[]interface{}{expensesToReturn, nil},
		1)

	actualExpenses, err := suite.service.SearchInPeriod(searchInPeriodCommand)

	require.NoError(suite.T(), err)
	for i, expectdExpense := range expensesToReturn {
		assertEqualsExpense(suite.T(), expectdExpense, actualExpenses[i])
	}
}

func (suite *ExpenseServiceTestSuite) TestGivenThatRepositoryFails_WhenSearchInPeriod_ThenReturnError() {
	searchInPeriodCommand, _ := NewSearchInPeriodCommand(
		time.Date(2022, 5, 23, 0, 0, 0, 0, time.Local),
		time.Date(2022, 8, 23, 0, 0, 0, 0, time.Local))

	suite.expenseRepositoryMock.MockSearchInPeriod(
		[]interface{}{searchInPeriodCommand.startDate, searchInPeriodCommand.endDate},
		[]interface{}{nil, errors.New("fail to get expenses")},
		1)

	actualExpenses, err := suite.service.SearchInPeriod(searchInPeriodCommand)

	require.ErrorAs(suite.T(), err, &UnexpectedError{})
	require.Nil(suite.T(), actualExpenses)
}

func (suite *ExpenseServiceTestSuite) getExpenses() []*models.Expense {
	return []*models.Expense{
		models.NewExpense(10.3, time.Date(2022, 5, 28, 0, 0, 0, 0, time.Local), "Lomitos", models.NewExpenseType("Servicios")),
		models.NewExpense(10.3, time.Date(2022, 7, 28, 0, 0, 0, 0, time.Local), "Lomitos", models.NewExpenseType("Servicios")),
	}
}

func getExpense() *models.Expense {
	return &models.Expense{
		Id:          pkg.NewUUID(),
		Amount:      100.50,
		ExpenseDate: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		Description: "Lomitos",
		ExpenseType: getExpenseType(),
	}

}

func getExpenseType() *models.ExpenseType {
	return &models.ExpenseType{
		Id:   pkg.NewUUID(),
		Name: "Delivery",
	}
}

func assertEqualsExpense(t *testing.T, expected *models.Expense, actual *models.Expense) {
	assert.Equal(t, expected.Id, actual.Id, "id are not equals")
	assert.Equal(t, expected.ExpenseType, actual.ExpenseType, "expenseTypes are not equals")
	assert.Equal(t, expected.Amount, actual.Amount, "amounts are not equals")
	assert.Equalf(t, expected.ExpenseDate, actual.ExpenseDate, "expenseDates are not equals")
	assert.Equalf(t, expected.Description, actual.Description, "descriptions are not equals")
}

func buildAddCommandFromExpense(expense *models.Expense) *AddCommand {
	return &AddCommand{
		amount:      expense.Amount,
		expenseDate: expense.ExpenseDate,
		description: expense.Description,
		expenseType: expense.ExpenseType,
	}
}
