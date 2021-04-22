package vacancy

import "post/internal/app/models"

type Repository interface {
	Create(vacancy models.Vacancy) (uint64, error)
	FindByID(id uint64) (*models.Vacancy, error)
}