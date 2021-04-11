package usecase

import "FL_2/model"

type ResponseUseCase interface {
	Create(response model.Response) (*model.Response, error)
	FindByOrderID(orderID uint64) ([]model.Response, error)
	Change(response model.Response) (*model.Response, error)
	Delete(response model.Response) error
}
