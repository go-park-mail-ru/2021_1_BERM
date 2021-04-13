package implementation

import (
	"FL_2/model"
	"FL_2/store"
	"github.com/pkg/errors"
)

const (
	vacancyUseCaseError = "Vacancy use case error"
)

type VacancyUseCase struct {
	store      store.Store
	mediaStore store.MediaStore
}

func (v *VacancyUseCase) Create(vacancy model.Vacancy) (*model.Vacancy, error) {
	id, err := v.store.Vacancy().Create(vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.ID = id
	err = v.supplementingTheVacancyModel(&vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	return &vacancy, err
}

func (v *VacancyUseCase) FindByID(id uint64) (*model.Vacancy, error) {
	vacancy, err := v.store.Vacancy().FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	err = v.supplementingTheVacancyModel(vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	return vacancy, nil
}

func (v *VacancyUseCase) supplementingTheVacancyModel(vacancy *model.Vacancy) error {
	u, err := v.store.User().FindByID(vacancy.UserId)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Login = u.Login
	image, err := v.mediaStore.Image().GetImage(u.Img)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Img = string(image)
	return nil
}
