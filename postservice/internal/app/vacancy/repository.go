package vacancy

import "post/internal/app/models"

type Repository interface {
	Create(vacancy models.Vacancy) (uint64, error)
	FindByID(id uint64) (*models.Vacancy, error)
	Change(vacancy models.Vacancy) error
	DeleteVacancy(id uint64) error
	FindByExecutorID(executorID uint64) ([]models.Vacancy, error)
	FindByCustomerID(customerID uint64) ([]models.Vacancy, error)
	UpdateExecutor(vacancy models.Vacancy) error
}