package vacancy

import "ff/internal/app/models"

type VacancyRepository interface {
	Create(vacancy models.Vacancy) (uint64, error)
	FindByID(id uint64) (*models.Vacancy, error)
}

