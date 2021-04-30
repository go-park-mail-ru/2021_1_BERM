package usecase

import (
	"context"
	"github.com/pkg/errors"
	"post/api"
	"post/internal/app/models"
	responseRepo "post/internal/app/response/repository"
)

const (
	responseUseCaseError = "Responce use case error"
)

type UseCase struct {
	ResponseRepo responseRepo.Repository
	UserRepo     api.UserClient
}

func NewUseCase(responseRepo responseRepo.Repository, userRepo api.UserClient) *UseCase {
	return &UseCase{
		ResponseRepo: responseRepo,
		UserRepo:     userRepo,
	}
}

func (u *UseCase) Create(response models.Response) (*models.Response, error) {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: response.UserID})
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}

	response.UserLogin = userR.GetLogin()
	response.UserImg = userR.GetImg()
	id, err := u.ResponseRepo.Create(response)
	response.ID = id
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	return &response, nil
}

func (u *UseCase) FindByPostID(postID uint64, orderResponse bool, vacancyResponse bool) ([]models.Response, error) {
	var responses []models.Response
	var err error
	if orderResponse {
		responses, err = u.ResponseRepo.FindByOrderPostID(postID)
	}
	if vacancyResponse {
		responses, err = u.ResponseRepo.FindByVacancyPostID(postID)
	}
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	for i, response := range responses {
		err := u.supplementingTheResponseModel(&response)
		if err != nil {
			return nil, errors.Wrap(err, responseUseCaseError)
		}
		responses[i].UserImg = response.UserImg
		responses[i].UserLogin = response.UserLogin
	}
	if responses == nil {
		return []models.Response{}, nil
	}
	return responses, nil
}

func (u *UseCase) Change(response models.Response) (*models.Response, error) {
	changedResponse := &models.Response{}
	var err error
	if response.OrderResponse {
		changedResponse, err = u.ResponseRepo.ChangeOrderResponse(response)
	}
	if response.VacancyResponse {
		changedResponse, err = u.ResponseRepo.ChangeVacancyResponse(response)
	}
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	err = u.supplementingTheResponseModel(changedResponse)
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	return changedResponse, nil
}

func (u *UseCase) Delete(response models.Response) error {
	var err error
	if response.OrderResponse {
		err = u.ResponseRepo.DeleteOrderResponse(response)
	}
	if response.VacancyResponse {
		err = u.ResponseRepo.DeleteVacancyResponse(response)
	}
	if err != nil {
		return errors.Wrap(err, responseUseCaseError)
	}
	return nil
}

func (u *UseCase) supplementingTheResponseModel(response *models.Response) error {
	user, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: response.UserID})
	if err != nil {
		return errors.Wrap(err, responseUseCaseError)
	}
	response.UserLogin = user.GetLogin()
	response.UserImg = user.GetImg()
	return nil
}
