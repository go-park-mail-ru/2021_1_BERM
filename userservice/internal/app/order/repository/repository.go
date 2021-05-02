package repository

import (
	"context"
	"user/internal/app/models"
)

type Repository interface {
	GetByID(ID uint64, ctx context.Context)  (models.OrderInfo, error)
}
