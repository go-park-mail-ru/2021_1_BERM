package order

import (
	"post/internal/app/models"
)

type Repository interface {
	Create(order models.Order) (uint64, error)
	Change(order models.Order) error
	DeleteOrder(id uint64) error
	FindByID(id uint64) (*models.Order, error)
	FindByExecutorID(executorID uint64) ([]models.Order, error)
	FindByCustomerID(customerID uint64) ([]models.Order, error)
	GetActualOrders() ([]models.Order, error)
	UpdateExecutor(order models.Order) error
}
