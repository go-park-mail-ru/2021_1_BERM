package response

import (
	"context"
	"github.com/pkg/errors"
	"post/api"
	"post/internal/app/models"
	responseRepo "post/internal/app/response"
	orderRepo "post/internal/app/order"
)

const (
	responseUseCaseError = "Responce use case error"
)

type UseCase struct {
	ResponseRepo responseRepo.Repository
	UserRepo     api.UserClient
	OrderRepo    orderRepo.Repository
}

func NewUseCase(responseRepo responseRepo.Repository,
	userRepo api.UserClient, orderR orderRepo.Repository) *UseCase {
	return &UseCase{
		ResponseRepo: responseRepo,
		UserRepo:     userRepo,
		OrderRepo: orderR,
	}
}

func (u *UseCase) Create(response models.Response, ctx context.Context) (*models.Response, error) {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: response.UserID})
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}

	response.UserLogin = userR.GetLogin()
	response.UserImg = userR.GetImg()
	id, err := u.ResponseRepo.Create(response, ctx)
	response.ID = id
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	return &response, nil
}

func (u *UseCase) FindByPostID(
	postID uint64,
	orderResponse bool,
	vacancyResponse bool,
	ctx context.Context) ([]models.Response, error) {
	var responses []models.Response
	var err error
	if orderResponse {
		responses, err = u.ResponseRepo.FindByOrderPostID(postID, ctx)
	}
	if vacancyResponse {
		responses, err = u.ResponseRepo.FindByVacancyPostID(postID, ctx)
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

func (u *UseCase) Change(response models.Response, ctx context.Context) (*models.Response, error) {
	changedResponse := &models.Response{}
	var err error
	if response.OrderResponse {
		changedResponse, err = u.ResponseRepo.ChangeOrderResponse(response, ctx)
	}
	if response.VacancyResponse {
		changedResponse, err = u.ResponseRepo.ChangeVacancyResponse(response, ctx)
	}
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	err = u.supplementingTheResponseModel(changedResponse)
	if err != nil {
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	o, err := u.OrderRepo.FindByID(response.PostID, ctx)
	if err != nil{
		return nil, errors.Wrap(err, responseUseCaseError)
	}
	if o.ID == response.UserID{
		return nil, errors.New(responseUseCaseError)
	}
	return changedResponse, nil



}

func (u *UseCase) Delete(response models.Response, ctx context.Context) error {
	var err error
	if response.OrderResponse {
		err = u.ResponseRepo.DeleteOrderResponse(response, ctx)
	}
	if response.VacancyResponse {
		err = u.ResponseRepo.DeleteVacancyResponse(response, ctx)
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
