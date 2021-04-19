package order

import "ff/internal/app/models"

type OrderUseCase interface {
	Create(order models.Order) (*models.Order, error)
	FindByID(id uint64) (*models.Order, error)
	FindByExecutorID(executorID uint64) ([]models.Order, error)
	FindByCustomerID(customerID uint64) ([]models.Order, error)
	GetActualOrders() ([]models.Order, error)
	SelectExecutor(order models.Order) error
	DeleteExecutor(order models.Order) error
}
