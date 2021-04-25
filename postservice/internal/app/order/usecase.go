package order

import (
	"post/internal/app/models"
)

type UseCase interface {
	Create(order models.Order) (*models.Order, error)
	FindByID(id uint64) (*models.Order, error)
	FindByUserID(userID uint64) ([]models.Order, error)
	GetActualOrders() ([]models.Order, error)
	SelectExecutor(order models.Order) error
	DeleteExecutor(order models.Order) error
}