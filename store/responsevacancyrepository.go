package store

import "FL_2/model"

type ResponseVacancyRepository interface {
	Create(response model.ResponseVacancy) (uint64, error)
	FindByVacancyID(id uint64) ([]model.ResponseVacancy, error)
	Change(response model.ResponseVacancy) (*model.ResponseVacancy, error)
	Delete(response model.ResponseVacancy) error
}