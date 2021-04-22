package usecase

import (
	"github.com/pkg/errors"
	"post/internal/app/models"
	responseRepo "post/internal/app/response/repository"
)

const (
	responseUseCaseError = "Responce use case error"
)

type UseCase struct {
	repo responseRepo.Repository
}

func NewUseCase(repo responseRepo.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Create(response models.Response) (*models.Response, error) {
	//TODO: grpc запрос за юзером
	user, err := u.repo.FindUserByID(response.UserID)
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	//TODO: grpc запрос за юзером
	user.Specializes, err = r.store.User().FindSpecializesByUserID(response.UserID)
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	response.UserLogin = user.Login
	response.UserImg = user.Img
	id, err := u.repo.Create(response)
	response.ID = id
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	response.ID = id
	//TODO: grpc-запрос за имгой
	img, err := r.mediaStore.Image().GetImage(response.UserImg)
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	response.UserImg = string(img)
	return &response, nil
}

func (u *UseCase) FindByPostID(postID uint64) ([]models.Response, error) {
	responses, err := u.repo.FindByPostID(postID)
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
	changedResponse, err := u.repo.Change(response)
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
	err := u.repo.Delete(response)
	if err != nil {
		return errors.Wrap(err, responseUseCaseError)
	}
	return nil
}

func (u *UseCase) supplementingTheResponseModel(response *models.Response) error {
	//TODO: grpc-запрос в юзер-репо
	user, err := o.store.User().FindUserByID(response.UserID)
	if err != nil {
		return errors.Wrap(err, responseUseCaseError)
	}
	//TODO: grpc-запрос в юзер-репо
	user.Specializes, err = o.store.User().FindSpecializesByUserID(response.UserID)
	if err != nil {
		return errors.Wrap(err, responseUseCaseError)
	}
	response.UserLogin = user.Login
	//TODO: grpc-запрос в картинки
	image, err := o.mediaStore.Image().GetImage(user.Img)
	if err != nil {
		return errors.Wrap(err, responseUseCaseError)
	}
	response.UserImg = string(image)
	return nil
}