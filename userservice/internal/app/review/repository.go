package review

import (
	"context"
	"user/internal/app/models"
)
//go:generate mockgen -destination=./mock/mock_repository.go -package=mock user/internal/app/review Repository
type Repository interface {
	Create(review models.Review, ctx context.Context) (*models.Review, error)
	GetAll(userId uint64, ctx context.Context) ([]models.Review, error)
	GetAvgScoreByUserId(userId uint64, ctx context.Context) (*models.UserReviewInfo, error)
}
