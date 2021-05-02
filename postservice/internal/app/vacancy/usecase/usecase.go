package usecase

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"post/api"
	"post/internal/app/models"
	vacancyRepo "post/internal/app/vacancy/repository"
	customErr "post/pkg/error"
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

func (u *UseCase) Create(vacancy models.Vacancy, ctx context.Context) (*models.Vacancy, error) {
	u.sanitizeVacancy(&vacancy)
	id, err := u.VacancyRepo.Create(vacancy, ctx)
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

func (u *UseCase) FindByID(id uint64, ctx context.Context) (*models.Vacancy, error) {
	vacancy, err := u.VacancyRepo.FindByID(id, ctx)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	err = u.supplementingTheVacancyModel(vacancy)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	return vacancy, nil
}

func (u *UseCase) GetActualVacancies(ctx context.Context) ([]models.Vacancy, error) {
	vacancies, err := u.VacancyRepo.GetActualVacancies(ctx)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	for i, vacancy := range vacancies {
		err = u.supplementingTheVacancyModel(&vacancy)
		if err != nil {
			return nil, errors.Wrap(err, vacancyUseCaseError)
		}
		vacancies[i] = vacancy
	}
	if vacancies == nil {
		return []models.Vacancy{}, nil
	}
	return vacancies, err
}

func (u *UseCase) ChangeVacancy(vacancy models.Vacancy, ctx context.Context) (models.Vacancy, error) {
	oldVacancy, err := u.VacancyRepo.FindByID(vacancy.ID, ctx)
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
	err = u.VacancyRepo.Change(vacancy, ctx)
	if err != nil {
		return models.Vacancy{}, errors.Wrap(err, vacancyUseCaseError)
	}
	if err := u.supplementingTheVacancyModel(&vacancy); err != nil {
		return models.Vacancy{}, errors.Wrap(err, vacancyUseCaseError)
	}
	return vacancy, nil
}

func (u *UseCase) DeleteVacancy(id uint64, ctx context.Context) error {
	err := u.VacancyRepo.DeleteVacancy(id, ctx)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) FindByUserID(userID uint64, ctx context.Context) ([]models.Vacancy, error) {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: userID})
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	isExecutor := userR.GetExecutor()
	var vacancies []models.Vacancy
	if isExecutor {
		vacancies, err = u.VacancyRepo.FindByExecutorID(userID, ctx)
	} else {
		vacancies, err = u.VacancyRepo.FindByCustomerID(userID, ctx)
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

func (u *UseCase) SelectExecutor(vacancy models.Vacancy, ctx context.Context) error {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: vacancy.ExecutorID})
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	//TODO: изменить интернал на экстернал
	if userR.GetExecutor() == false {
		return customErr.ErrorUserNotExecutor
	}
	//TODO: изменить интернал на экстернал
	if vacancy.ExecutorID == vacancy.CustomerID {
		return customErr.ErrorSameID
	}
	err = u.VacancyRepo.UpdateExecutor(vacancy, ctx)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) DeleteExecutor(vacancy models.Vacancy, ctx context.Context) error {
	vacancy.ExecutorID = 0
	err := u.VacancyRepo.UpdateExecutor(vacancy, ctx)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) CloseVacancy(vacancyID uint64, ctx context.Context) error {
	vacancy, err := u.VacancyRepo.FindByID(vacancyID, ctx)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	err = u.VacancyRepo.DeleteVacancy(vacancyID, ctx)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	_, err = u.VacancyRepo.CreateArchive(*vacancy, ctx)
	if err != nil {
		return errors.Wrap(err, vacancyUseCaseError)
	}
	return nil
}

func (u *UseCase) GetArchiveVacancies(ctx context.Context) ([]models.Vacancy, error) {
	vacancies, err := u.VacancyRepo.GetArchiveVacancies(ctx)
	if err != nil {
		return nil, errors.Wrap(err, vacancyUseCaseError)
	}
	for i, vacancy := range vacancies {
		err = u.supplementingTheVacancyModel(&vacancy)
		if err != nil {
			return nil, errors.Wrap(err, vacancyUseCaseError)
		}
		vacancies[i] = vacancy
	}
	if vacancies == nil {
		return []models.Vacancy{}, nil
	}
	return vacancies, err
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
