package response

import (
	"post/internal/app/models"
)

type Repository interface {
	Create(response models.Response) (uint64, error)
	FindByPostID(id uint64) ([]models.Response, error)
	Change(response models.Response) (*models.Response, error)
	Delete(response models.Response) error
}
