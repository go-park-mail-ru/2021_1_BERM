package vacancy

import (
	"post/internal/app/models"
)

type UseCase interface {
	Create(vacancy models.Vacancy) (*models.Vacancy, error)
	FindByID(id uint64) (*models.Vacancy, error)
	ChangeVacancy(vacancy models.Vacancy) (models.Vacancy, error)
	DeleteVacancy(id uint64) error
	FindByUserID(userID uint64) ([]models.Vacancy, error)
	SelectExecutor(vacancy models.Vacancy) error
	DeleteExecutor(vacancy models.Vacancy) error
}
