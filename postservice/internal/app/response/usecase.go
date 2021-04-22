package response

import (
	"post/internal/app/models"
)

type UseCase interface {
	Create(response models.Response) (*models.Response, error)
	FindByPostID(orderID uint64) ([]models.Response, error)
	Change(response models.Response) (*models.Response, error)
	Delete(response models.Response) error
}
