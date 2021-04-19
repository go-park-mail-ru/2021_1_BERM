package vacancy

import "ff/internal/app/models"

type VacancyUseCase interface {
	Create(vacancy models.Vacancy) (*models.Vacancy, error)
	FindByID(id uint64) (*models.Vacancy, error)
}
