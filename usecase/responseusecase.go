package usecase

import "FL_2/model"

type ResponseUseCase interface {
	Create(response model.Response) (*model.Response, error)
	FindByID(id uint64) ([]model.Response, error)
}
