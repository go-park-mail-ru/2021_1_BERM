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
	UserRepo api.UserClient
}

func NewUseCase(responseRepo responseRepo.Repository, userRepo api.UserClient) *UseCase {
	return &UseCase{
		ResponseRepo: responseRepo,
		UserRepo: userRepo,
	}
}

func (u *UseCase) Create(response models.Response) (*models.Response, error) {
	//TODO: grpc запрос за юзером
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

func (u *UseCase) FindByPostID(postID uint64) ([]models.Response, error) {
	responses, err := u.ResponseRepo.FindByPostID(postID)
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
	changedResponse, err := u.ResponseRepo.Change(response)
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
	err := u.ResponseRepo.Delete(response)
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
