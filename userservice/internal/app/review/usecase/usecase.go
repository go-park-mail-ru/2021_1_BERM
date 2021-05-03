package usecase

import (
	"context"
	"user/internal/app/models"
)

type UseCase interface {
	Create(review models.Review, ctx context.Context) (*models.Review, error)
	GetAllReviewByUserId(userId uint64, ctx context.Context) (*models.UserReviews, error)
}

