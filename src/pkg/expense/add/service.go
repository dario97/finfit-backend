package add

const invalidExpenseTypeErrorMsg = "the expenseDBModel type doesn't exists"

type Service interface {
	Add(command command) (*createdExpense, error)
}

type Repository interface {
	AddExpense(entity expense) (*createdExpense, error)
	GetExpenseTypeById(id uint64) (*expenseType, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Add(command command) (*createdExpense, error) {
	expenseType, repositoryError := s.repository.GetExpenseTypeById(command.expenseType.id)

	if repositoryError != nil {
		return nil, UnexpectedError{Msg: repositoryError.Error()}
	}

	if expenseType == nil {
		return nil, InvalidExpenseTypeError{Msg: invalidExpenseTypeErrorMsg}
	}

	expenseToCreate := expense{
		amount:      command.amount,
		expenseDate: command.expenseDate,
		description: command.description,
		expenseType: command.expenseType,
	}

	createdExpense, repoError := s.repository.AddExpense(expenseToCreate)

	if repoError != nil {
		return nil, UnexpectedError{Msg: repoError.Error()}
	}

	return createdExpense, nil
}
