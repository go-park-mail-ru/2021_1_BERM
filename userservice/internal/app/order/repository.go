package order

import (
	"context"
	"user/internal/app/models"
)
//go:generate mockgen -destination=./mock/mock_repository.go -package=mock user/internal/app/order Repository
type Repository interface {
	GetByID(ID uint64, ctx context.Context) (*models.OrderInfo, error)
}
