package response

import (
	"post/internal/app/models"
)

type Repository interface {
	Create(response models.Response) (uint64, error)
	FindByOrderPostID(id uint64) ([]models.Response, error)
	FindByVacancyPostID(id uint64) ([]models.Response, error)
	ChangeOrderResponse(response models.Response) (*models.Response, error)
	ChangeVacancyResponse(response models.Response) (*models.Response, error)
	DeleteOrderResponse(response models.Response) error
	DeleteVacancyResponse(response models.Response) error
}
