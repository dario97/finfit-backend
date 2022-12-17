package expensetype

import (
	"finfit-backend/internal/domain/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(id uuid.UUID) (*models.ExpenseType, error)
}
type Service interface {
	GetById(id uuid.UUID) (*models.ExpenseType, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s service) GetById(id uuid.UUID) (*models.ExpenseType, error) {
	expenseType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}

	return expenseType, nil
}

type UnexpectedError struct {
	Msg string
}

func (receiver UnexpectedError) Error() string {
	return receiver.Msg
}
