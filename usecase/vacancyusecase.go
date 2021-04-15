package usecase

import "FL_2/model"

type VacancyUseCase interface {
	Create(vacancy model.Vacancy) (*model.Vacancy, error)
	FindByID(id uint64) (*model.Vacancy, error)
}
