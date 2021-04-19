package order

import "ff/internal/app/models"

type OrderRepository interface {
	Create(order models.Order) (uint64, error)
	FindByID(id uint64) (*models.Order, error)
	FindByExecutorID(executorID uint64) ([]models.Order, error)
	FindByCustomerID(customerID uint64) ([]models.Order, error)
	GetActualOrders() ([]models.Order, error)
	UpdateExecutor(order models.Order) error
}
