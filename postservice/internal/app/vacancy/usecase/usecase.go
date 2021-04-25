package usecase

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"post/api"
	"post/internal/app/models"
	vacancyRepo "post/internal/app/vacancy/repository"
)

const (
	vacancyUseCaseError = "Vacancy use case error"
)


type UseCase struct {
	VacancyRepo vacancyRepo.Repository
	UserRepo api.UserClient
}

func NewUseCase(vacancyRepo vacancyRepo.Repository, userRepo api.UserClient) *UseCase {
	return &UseCase{
		VacancyRepo: vacancyRepo,
		UserRepo: userRepo,
	}
}

func (u *UseCase) Create(vacancy models.Vacancy) (*models.Vacancy, error) {
	u.sanitizeVacancy(&vacancy)
	id, err := u.VacancyRepo.Create(vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.ID = id
	err = u.supplementingTheVacancyModel(&vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	return &vacancy, err
}

func (u *UseCase) FindByID(id uint64) (*models.Vacancy, error) {
	vacancy, err := u.VacancyRepo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	err = u.supplementingTheVacancyModel(vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	return vacancy, nil
}

func (u *UseCase) supplementingTheVacancyModel(vacancy *models.Vacancy) error {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: vacancy.UserID})
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Login = userR.GetLogin()
	vacancy.Img = userR.GetImg()
	return nil
}

func (u *UseCase) sanitizeVacancy(vacancy *models.Vacancy) {
	sanitizer := bluemonday.UGCPolicy()
	vacancy.VacancyName = sanitizer.Sanitize(vacancy.VacancyName)
	vacancy.Description = sanitizer.Sanitize(vacancy.Description)
	vacancy.Category = sanitizer.Sanitize(vacancy.Category)
}
