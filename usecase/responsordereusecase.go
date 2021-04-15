package usecase

import "FL_2/model"

type ResponseOrderUseCase interface {
	Create(response model.ResponseOrder) (*model.ResponseOrder, error)
	FindByVacancyID(orderID uint64) ([]model.ResponseOrder, error)
	Change(response model.ResponseOrder) (*model.ResponseOrder, error)
	Delete(response model.ResponseOrder) error
}
