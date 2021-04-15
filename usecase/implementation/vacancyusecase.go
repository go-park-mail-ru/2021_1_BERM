package implementation

import (
	"FL_2/model"
	"FL_2/store"
	"github.com/microcosm-cc/bluemonday"
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
	v.sanitizeVacancy(&vacancy)
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
	user, err := v.store.User().FindUserByID(vacancy.UserID)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Login = user.Login
	image, err := v.mediaStore.Image().GetImage(user.Img)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Img = string(image)
	return nil
}

func (v *VacancyUseCase) sanitizeVacancy(vacancy *model.Vacancy) {
	sanitizer := bluemonday.UGCPolicy()
	vacancy.VacancyName = sanitizer.Sanitize(vacancy.VacancyName)
	vacancy.Description = sanitizer.Sanitize(vacancy.Description)
	vacancy.Category = sanitizer.Sanitize(vacancy.Category)
}
