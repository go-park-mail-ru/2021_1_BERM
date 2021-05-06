package review

import (
	"context"
	"user/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock user/internal/app/review UseCase
type UseCase interface {
	Create(review models.Review, ctx context.Context) (*models.Review, error)
	GetAllReviewByUserId(userId uint64, ctx context.Context) (*models.UserReviews, error)
}

