package implementation

import (
	"FL_2/model"
	"FL_2/store"
)

type VacancyUseCase struct {
	store store.Store
}

func (v *VacancyUseCase)Create(vacancy model.Vacancy) (*model.Vacancy, error){
	id, err := v.store.Vacancy().Create(vacancy)
	if err != nil{
		return nil, err
	}
	vacancy.Id = id
	return &vacancy, err
}

func (v *VacancyUseCase)FindByID(id uint64) (*model.Vacancy, error){
	vacancy, err := v.store.Vacancy().FindByID(id)
	if err != nil{
		return nil, err
	}
	return vacancy, nil
}