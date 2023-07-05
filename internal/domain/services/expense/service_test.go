package expense_test

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expense"
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
	expenseRepositoryMock  *expense.RepositoryMock
	expenseTypeServiceMock *expensetype.ServiceMock
	service                expense.Service
}

func (suite *ExpenseServiceTestSuite) SetupSuite() {
	suite.expenseRepositoryMock = expense.NewRepositoryMock()
	suite.expenseTypeServiceMock = expensetype.NewServiceMock()
	suite.service = expense.NewService(suite.expenseRepositoryMock, suite.expenseTypeServiceMock)
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
	expenseToCreate := suite.getExpense1()
	expectedCreatedExpense := expenseToCreate

	suite.expenseRepositoryMock.MockAdd([]interface{}{expenseToCreate}, []interface{}{expectedCreatedExpense, nil}, 1)
	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType().Id()}, []interface{}{expenseToCreate.ExpenseType(), nil}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), err, "Error must to be nil")
	assertEqualsExpense(suite.T(), expectedCreatedExpense, actualCreatedExpense)
}

func (suite *ExpenseServiceTestSuite) TestGivenThatFailToGetExpenseType_WhenAdd_ThenReturnError() {
	expenseToCreate := suite.getExpense1()

	expenseTypeServiceError := errors.New("fail to get expense type")
	expectedError := expense.UnexpectedError{
		Msg: expenseTypeServiceError.Error(),
	}

	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType().Id()}, []interface{}{nil, expenseTypeServiceError}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), actualCreatedExpense)
	assert.NotNil(suite.T(), err, "Error must not be nil")
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *ExpenseServiceTestSuite) TestGivenThatExpenseTypeNotExists_WhenAdd_ThenReturnError() {
	expenseToCreate := suite.getExpense1()
	expectedError := expense.InvalidExpenseTypeError{
		Msg: "the expense type doesn't exists",
	}

	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType().Id()}, []interface{}{nil, nil}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), actualCreatedExpense)
	assert.NotNil(suite.T(), err, "Error must not be nil")
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *ExpenseServiceTestSuite) TestGivenThatSaveExpenseIntoDatabaseFails_WhenAdd_ThenReturnError() {
	expenseToCreate := suite.getExpense1()

	repoError := errors.New("fail to save expense")
	expectedError := expense.UnexpectedError{
		Msg: repoError.Error(),
	}

	suite.expenseTypeServiceMock.MockGetByID([]interface{}{expenseToCreate.ExpenseType().Id()}, []interface{}{expenseToCreate.ExpenseType(), nil}, 1)
	suite.expenseRepositoryMock.MockAdd([]interface{}{expenseToCreate}, []interface{}{nil, repoError}, 1)

	actualCreatedExpense, err := suite.service.Add(buildAddCommandFromExpense(expenseToCreate))

	assert.Nil(suite.T(), actualCreatedExpense)
	assert.NotNil(suite.T(), err, "Error must not be nil")
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *ExpenseServiceTestSuite) TestGivenAPeriod_WhenSearchInPeriod_ThenReturnAListOfExpenses() {
	expensesToReturn := suite.getExpenses()

	searchInPeriodCommand, _ := expense.NewSearchInPeriodCommand(
		time.Date(2022, 5, 23, 0, 0, 0, 0, time.Local),
		time.Date(2022, 8, 23, 0, 0, 0, 0, time.Local))

	suite.expenseRepositoryMock.MockSearchInPeriod(
		[]interface{}{searchInPeriodCommand.StartDate(), searchInPeriodCommand.EndDate()},
		[]interface{}{expensesToReturn, nil},
		1)

	actualExpenses, err := suite.service.SearchInPeriod(searchInPeriodCommand)

	require.NoError(suite.T(), err)
	for i, expectdExpense := range expensesToReturn {
		assertEqualsExpense(suite.T(), expectdExpense, actualExpenses[i])
	}
}

func (suite *ExpenseServiceTestSuite) TestGivenThatRepositoryFails_WhenSearchInPeriod_ThenReturnError() {
	searchInPeriodCommand, _ := expense.NewSearchInPeriodCommand(
		time.Date(2022, 5, 23, 0, 0, 0, 0, time.Local),
		time.Date(2022, 8, 23, 0, 0, 0, 0, time.Local))

	suite.expenseRepositoryMock.MockSearchInPeriod(
		[]interface{}{searchInPeriodCommand.StartDate(), searchInPeriodCommand.EndDate()},
		[]interface{}{nil, errors.New("fail to get expenses")},
		1)

	actualExpenses, err := suite.service.SearchInPeriod(searchInPeriodCommand)

	require.ErrorAs(suite.T(), err, &expense.UnexpectedError{})
	require.Nil(suite.T(), actualExpenses)
}

func (suite *ExpenseServiceTestSuite) TestGivenExpensesToAdd_WhenAddAll_ThenReturnCreatedExpenses() {
	expensesToAdd := suite.getExpenses()
	expectedAddedExpenses := expensesToAdd

	actualCreatedExpenses, err := suite.service.AddAll(suite.buildAddAllCommandFromExpenses(expensesToAdd))

	require.NoError(suite.T(), err, "Error must to be nil")
	suite.assertEqualsExpenses(expectedAddedExpenses, actualCreatedExpenses)
}

func (suite *ExpenseServiceTestSuite) TestGivenRepositoryFails_WhenAddAll_ThenReturnUnexpectedError() {
	actualCreatedExpenses, err := suite.service.AddAll(suite.buildAddAllCommandFromExpenses(suite.getExpenses()))

	require.ErrorAs(suite.T(), err, &expense.UnexpectedError{})
	require.Nil(suite.T(), actualCreatedExpenses)
}

func (suite *ExpenseServiceTestSuite) assertEqualsExpenses(expectedAddedExpenses []*models.Expense, actualCreatedExpenses []*models.Expense) {
	assert.Equal(suite.T(), len(expectedAddedExpenses), len(actualCreatedExpenses))
	for i, expectedExpense := range expectedAddedExpenses {
		assertEqualsExpense(suite.T(), expectedExpense, actualCreatedExpenses[i])
	}
}

func (suite *ExpenseServiceTestSuite) getExpenses() []*models.Expense {
	expense1 := suite.getExpense1()
	expense2 := suite.getExpense2()

	return []*models.Expense{
		expense1,
		expense2,
	}
}

func (suite *ExpenseServiceTestSuite) getExpense1() *models.Expense {
	newExpense, _ := models.NewExpense(suite.getMoney(), time.Date(2022, 5, 28, 0, 0, 0, 0, time.Local), "Lomitos", suite.getExpenseType())
	return newExpense
}

func (suite *ExpenseServiceTestSuite) getExpense2() *models.Expense {
	newExpense, _ := models.NewExpense(suite.getMoney(), time.Date(2022, 7, 28, 0, 0, 0, 0, time.Local), "Lomitos", suite.getExpenseType())
	return newExpense
}

func (suite *ExpenseServiceTestSuite) buildAddAllCommandFromExpenses(expenses []*models.Expense) *expense.AddAllCommand {
	addCommands := []*expense.AddCommand{}
	for _, domainExpense := range expenses {
		addCommands = append(addCommands, buildAddCommandFromExpense(domainExpense))
	}

	command, _ := expense.NewAddAllCommand(addCommands)
	return command
}

func (suite *ExpenseServiceTestSuite) getExpenseType() *models.ExpenseType {
	expenseType, _ := models.NewExpenseType("Delivery")
	return expenseType
}

func (suite *ExpenseServiceTestSuite) getMoney() *models.Money {
	money, _ := models.NewMoney(10.3, "ARS")
	return money
}

func assertEqualsExpense(t *testing.T, expected *models.Expense, actual *models.Expense) {
	assert.Equal(t, expected.Id(), actual.Id(), "id are not equals")
	assert.Equal(t, expected.ExpenseType(), actual.ExpenseType(), "expenseTypes are not equals")
	assert.Equal(t, expected.Amount(), actual.Amount(), "amounts are not equals")
	assert.Equalf(t, expected.ExpenseDate(), actual.ExpenseDate(), "expenseDates are not equals")
	assert.Equalf(t, expected.Description(), actual.Description(), "descriptions are not equals")
}

func buildAddCommandFromExpense(domainExpense *models.Expense) *expense.AddCommand {
	addCommand, _ := expense.NewAddCommand(domainExpense.Amount().Amount(), domainExpense.Amount().Currency(), domainExpense.ExpenseDate(), domainExpense.Description(), domainExpense.ExpenseType().Id())
	return addCommand
}
