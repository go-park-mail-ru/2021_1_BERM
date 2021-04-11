package usecase

import "FL_2/model"

type OrderUseCase interface {
	Create(order model.Order) (*model.Order, error)
	FindByID(id uint64) (*model.Order, error)
	FindByExecutorID(executorID uint64) ([]model.Order, error)
	FindByCustomerID(customerID uint64) ([]model.Order, error)
	GetActualOrders()([]model.Order, error)
}

