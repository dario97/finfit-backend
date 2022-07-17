package registry

import "finfit-backend/src/domain/use_cases/service"

type ServiceRegistry struct {
	repositoryRegistry RepositoryRegistry
}

func NewServiceRegistry(repositoryRegistry RepositoryRegistry) ServiceRegistry {
	return ServiceRegistry{
		repositoryRegistry: repositoryRegistry,
	}
}

func (receiver ServiceRegistry) GetExpenseService() service.ExpenseService {
	return service.NewExpenseService(receiver.repositoryRegistry.GetExpenseRepository(), receiver.GetExpenseTypeService())
}

func (receiver ServiceRegistry) GetExpenseTypeService() service.ExpenseTypeService {
	return service.NewExpenseTypeService(receiver.repositoryRegistry.GetExpenseTypeRepository())
}
