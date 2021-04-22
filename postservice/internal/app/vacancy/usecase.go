package vacancy

import (
	"post/internal/app/models"
)

type UseCase interface {
	Create(vacancy models.Vacancy) (*models.Vacancy, error)
	FindByID(id uint64) (*models.Vacancy, error)
}

