package store

import "FL_2/model"

//go:generate mockgen -destination=mock/mock_order_repo.go -package=mock FL_2/store OrderRepository
type OrderRepository interface {
	Create(order model.Order) (uint64, error)
	FindByID(id uint64) (*model.Order, error)
	FindByExecutorID(executorID uint64) ([]model.Order, error)
	FindByCustomerID(customerID uint64) ([]model.Order, error)
	GetActualOrders() ([]model.Order, error)
	UpdateExecutor(order model.Order) error
}
