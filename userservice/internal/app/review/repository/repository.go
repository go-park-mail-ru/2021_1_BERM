package repository

import (
	"context"
	"user/internal/app/models"
)

type Repository interface {
	Create(review models.Review, ctx context.Context) (*models.Review, error)
	GetAll(userId uint64, ctx context.Context) ([]models.Review, error)
	GetAvgScoreByUserId(userId uint64, ctx context.Context) (uint8, error)
}
