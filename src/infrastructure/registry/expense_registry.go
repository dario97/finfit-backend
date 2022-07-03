package registry

import (
	repository2 "finfit/finfit-backend/src/domain/repository"
	"finfit/finfit-backend/src/domain/use_cases/service"
	"finfit/finfit-backend/src/infrastructure/repository"
	"finfit/finfit-backend/src/interfaces/controller"
)

func (r expenseRegistry) GetExpenseController() controller.ExpenseController {
	return controller.NewExpenseController(r.GetExpenseService())
}

func (r expenseRegistry) GetExpenseService() service.ExpenseService {
	return service.NewExpenseService(r.GetExpenseRepository(), r.GetExpenseTypeService())
}

func (r expenseRegistry) GetExpenseRepository() repository2.ExpenseRepository {
	return repository.NewExpenseRepository(r.db)
}

func (r expenseRegistry) GetExpenseTypeService() service.ExpenseTypeService {
	return service.NewExpenseTypeService(r.GetExpenseTypeRepository())
}

func (r expenseRegistry) GetExpenseTypeRepository() repository2.ExpenseTypeRepository {
	return repository.NewExpenseTypeRepository(r.db)
}
