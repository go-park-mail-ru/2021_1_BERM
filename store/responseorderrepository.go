package store

import "FL_2/model"

//go:generate mockgen -destination=mock/mock_responseorder_repo.go -package=mock FL_2/store ResponseOrderRepository
type ResponseOrderRepository interface {
	Create(response model.ResponseOrder) (uint64, error)
	FindByOrderID(id uint64) ([]model.ResponseOrder, error)
	Change(response model.ResponseOrder) (*model.ResponseOrder, error)
	Delete(response model.ResponseOrder) error
}
