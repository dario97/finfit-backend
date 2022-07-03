package registry

import (
	repository2 "finfit/finfit-backend/domain/repository"
	expenseservice "finfit/finfit-backend/domain/use_cases/service"
	"finfit/finfit-backend/infrastructure/repository"
	"finfit/finfit-backend/interfaces/controller"
)

func (r expenseRegistry) NewExpenseController() controller.ExpenseController {
	return controller.NewExpenseController(r.NewExpenseService())
}

func (r expenseRegistry) NewExpenseService() expenseservice.ExpenseService {
	return expenseservice.NewExpenseService(r.NewExpenseRepository())
}

func (r expenseRegistry) NewExpenseRepository() repository2.ExpenseRepository {
	return repository.NewExpenseRepository(r.db)
}
