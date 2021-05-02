package response

import (
	"context"
	"post/internal/app/models"
)

type UseCase interface {
	Create(response models.Response, ctx context.Context) (*models.Response, error)
	FindByPostID(postID uint64, orderResponse bool, vacancyResponse bool, ctx context.Context) ([]models.Response, error)
	Change(response models.Response, ctx context.Context) (*models.Response, error)
	Delete(response models.Response, ctx context.Context) error
}
