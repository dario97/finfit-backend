package registry

import (
	"finfit-backend/src/infrastructure/interfaces/controller"
)

type ControllerRegistry struct {
	serviceRegistry ServiceRegistry
}

func NewControllerRegistry(serviceRegistry ServiceRegistry) ControllerRegistry {
	return ControllerRegistry{
		serviceRegistry: serviceRegistry,
	}
}

func (receiver ControllerRegistry) GetExpenseController() controller.ExpenseController {
	return controller.NewExpenseController(receiver.serviceRegistry.GetExpenseService(), GetGenericFieldsValidator())
}
