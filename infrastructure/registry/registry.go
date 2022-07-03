package registry

import (
	"finfit/finfit-backend/interfaces/controller"
	"github.com/jinzhu/gorm"
)

type Registry interface {
	NewAppController() controller.AppController
}

type expenseRegistry struct {
	db *gorm.DB
}

func NewRegistry(db *gorm.DB) Registry {
	return &expenseRegistry{db: db}
}

func (r *expenseRegistry) NewAppController() controller.AppController {
	return r.NewAppController()
}
