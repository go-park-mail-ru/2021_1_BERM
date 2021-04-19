package usecase

import (
	"ff/internal/app/image"
	"ff/internal/app/models"
	"ff/internal/app/order"
	"ff/internal/app/user"
	"ff/internal/app/vacancy"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
)

const (
	vacancyUseCaseError = "Vacancy use case errors"
)

type VacancyUseCase struct {
	userRep    user.UserRepository
	orderRep   order.OrderRepository
	vacancyRep vacancy.VacancyRepository
	image      image.ImageRepository
}

func (v *VacancyUseCase) Create(vacancy models.Vacancy) (*models.Vacancy, error) {
	v.sanitizeVacancy(&vacancy)
	id, err := v.vacancyRep.Create(vacancy)
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

func (v *VacancyUseCase) FindByID(id uint64) (*models.Vacancy, error) {
	vacancy, err := v.vacancyRep.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	err = v.supplementingTheVacancyModel(vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	return vacancy, nil
}

func (v *VacancyUseCase) supplementingTheVacancyModel(vacancy *models.Vacancy) error {
	user, err := v.userRep.FindUserByID(vacancy.UserID)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Login = user.Login
	image, err := v.image.GetImage(user.Img)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Img = string(image)
	return nil
}

func (v *VacancyUseCase) sanitizeVacancy(vacancy *models.Vacancy) {
	sanitizer := bluemonday.UGCPolicy()
	vacancy.VacancyName = sanitizer.Sanitize(vacancy.VacancyName)
	vacancy.Description = sanitizer.Sanitize(vacancy.Description)
	vacancy.Category = sanitizer.Sanitize(vacancy.Category)
}
