package usecase

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"post/internal/app/models"
	vacancyRepo "post/internal/app/vacancy/repository"
)

const (
	vacancyUseCaseError = "Vacancy use case error"
)


type UseCase struct {
	repo vacancyRepo.Repository
}

func NewUseCase(repo vacancyRepo.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Create(vacancy models.Vacancy) (*models.Vacancy, error) {
	u.sanitizeVacancy(&vacancy)
	id, err := u.repo.Create(vacancy)
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
	vacancy, err := u.repo.FindByID(id)
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
	//TODO: grpc-запрос за юзером
	user, err := v.store.User().FindUserByID(vacancy.UserID)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Login = user.Login
	//TODO: grpc-запрос за img
	image, err := v.mediaStore.Image().GetImage(user.Img)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	vacancy.Img = string(image)
	return nil
}

func (u *UseCase) sanitizeVacancy(vacancy *models.Vacancy) {
	sanitizer := bluemonday.UGCPolicy()
	vacancy.VacancyName = sanitizer.Sanitize(vacancy.VacancyName)
	vacancy.Description = sanitizer.Sanitize(vacancy.Description)
	vacancy.Category = sanitizer.Sanitize(vacancy.Category)
}
