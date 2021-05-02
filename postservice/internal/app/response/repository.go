package response

import (
	"context"
	"post/internal/app/models"
)

type Repository interface {
	Create(response models.Response, ctx context.Context) (uint64, error)
	FindByOrderPostID(id uint64, ctx context.Context) ([]models.Response, error)
	FindByVacancyPostID(id uint64, ctx context.Context) ([]models.Response, error)
	ChangeOrderResponse(response models.Response, ctx context.Context) (*models.Response, error)
	ChangeVacancyResponse(response models.Response, ctx context.Context) (*models.Response, error)
	DeleteOrderResponse(response models.Response, ctx context.Context) error
	DeleteVacancyResponse(response models.Response, ctx context.Context) error
}
