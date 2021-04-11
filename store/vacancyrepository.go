package store

import "FL_2/model"

type VacancyRepository interface {
	Create(vacancy model.Vacancy) (uint64, error)
	FindByID(id uint64) (*model.Vacancy, error)
}
