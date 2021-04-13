package usecase

import "FL_2/model"

type ResponseVacancyUseCase interface {
	Create(response model.ResponseVacancy) (*model.ResponseVacancy, error)
	FindByVacancyID(orderID uint64) ([]model.ResponseVacancy, error)
	Change(response model.ResponseVacancy) (*model.ResponseVacancy, error)
	Delete(response model.ResponseVacancy) error
}
