package usecase

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"post/api"
	"post/internal/app/models"
	vacancyRepo "post/internal/app/vacancy/repository"
	"post/pkg/Error"
)

const (
	vacancyUseCaseError = "Vacancy use case error"
)

type UseCase struct {
	VacancyRepo vacancyRepo.Repository
	UserRepo    api.UserClient
}

func NewUseCase(vacancyRepo vacancyRepo.Repository, userRepo api.UserClient) *UseCase {
	return &UseCase{
		VacancyRepo: vacancyRepo,
		UserRepo:    userRepo,
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

func (u *UseCase) ChangeVacancy(vacancy models.Vacancy) (models.Vacancy, error) {
	oldVacancy, err := u.VacancyRepo.FindByID(vacancy.ID)
	if err != nil {
		return models.Vacancy{}, errors.Wrap(err, vacancyUseCaseError)
	}
	u.sanitizeVacancy(&vacancy)
	if vacancy.VacancyName == "" {
		vacancy.VacancyName = oldVacancy.VacancyName
	}
	if vacancy.Category == "" {
		vacancy.Category = oldVacancy.Category
	}
	if vacancy.Description == "" {
		vacancy.Description = oldVacancy.Description
	}
	if vacancy.Salary == 0 {
		vacancy.Salary = oldVacancy.Salary
	}
	vacancy.CustomerID = oldVacancy.CustomerID
	vacancy.ExecutorID = oldVacancy.ExecutorID
	err = u.VacancyRepo.Change(vacancy)
	if err != nil {
		return models.Vacancy{}, errors.Wrap(err, vacancyUseCaseError)
	}
	if err := u.supplementingTheVacancyModel(&vacancy); err != nil {
		return models.Vacancy{}, errors.Wrap(err, vacancyUseCaseError)
	}
	return vacancy, nil
}

func (u *UseCase) DeleteVacancy(id uint64) error {
	err := u.VacancyRepo.DeleteVacancy(id)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) FindByUserID(userID uint64) ([]models.Vacancy, error) {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: userID})
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	isExecutor := userR.GetExecutor()
	var vacancies []models.Vacancy
	if isExecutor {
		vacancies, err = u.VacancyRepo.FindByExecutorID(userID)
	} else {
		vacancies, err = u.VacancyRepo.FindByCustomerID(userID)
	}
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	for _, vacancy := range vacancies {
		err = u.supplementingTheVacancyModel(&vacancy)
		if err != nil {
			return nil, errors.Wrap(err, vacancyUseCaseError)
		}
	}
	if vacancies == nil {
		return []models.Vacancy{}, nil
	}
	return vacancies, nil
}

func (u *UseCase) SelectExecutor(vacancy models.Vacancy) error {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: vacancy.ExecutorID})
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	//TODO: изменить интернал на экстернал
	if userR.GetExecutor() == false {
		return &Error.Error{
			Err: errors.New("Select user not executor"),
			ErrorDescription: map[string]interface{}{
				"Error": Error.InternalServerErrorDescription["Error"]},
			InternalError: true,
		}
	}
	//TODO: изменить интернал на экстернал
	if vacancy.ExecutorID == vacancy.CustomerID {
		return &Error.Error{
			Err: errors.New("Executor and customer ID are the same"),
			ErrorDescription: map[string]interface{}{
				"Error": Error.InternalServerErrorDescription["Error"]},
			InternalError: true,
		}
	}
	err = u.VacancyRepo.UpdateExecutor(vacancy)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) DeleteExecutor(vacancy models.Vacancy) error {
	vacancy.ExecutorID = 0
	err := u.VacancyRepo.UpdateExecutor(vacancy)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) supplementingTheVacancyModel(vacancy *models.Vacancy) error {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: vacancy.CustomerID})
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
