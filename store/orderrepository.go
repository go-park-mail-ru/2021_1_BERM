package store

import "FL_2/model"

type OrderRepository interface {
	Create(order model.Order) (uint64, error)
	FindByID(id uint64) (*model.Order, error)
	FindByExecutorID(executorID uint64) ([]model.Order, error)
	FindByCustomerID(customerID uint64) ([]model.Order, error)
	GetActualOrders() ([]model.Order, error)
	UpdateExecutor(order model.Order) error
}
