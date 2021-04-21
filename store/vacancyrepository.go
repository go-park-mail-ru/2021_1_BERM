package store

import "FL_2/model"

//go:generate mockgen -destination=mock/mock_vacansy_repo.go -package=mock FL_2/store VacancyRepository
type VacancyRepository interface {
	Create(vacancy model.Vacancy) (uint64, error)
	FindByID(id uint64) (*model.Vacancy, error)
}
