package implementation

import (
	"FL_2/model"
	"FL_2/store"
	"github.com/pkg/errors"
)

const (
	responceUseCaseError = "Responce use case error"
)

type ResponseOrderUseCase struct {
	store      store.Store
	mediaStore store.MediaStore
}

func (r *ResponseOrderUseCase) Create(response model.ResponseOrder) (*model.ResponseOrder, error) {
	user, err := r.store.User().FindByID(response.UserID)
	if err != nil {
		return nil, errors.Wrap(err, responceUseCaseError)
	}
	response.UserLogin = user.Login
	response.UserImg = user.Img
	id, err := r.store.ResponseOrder().Create(response)
	response.ID = id
	if err != nil {
		return nil, errors.Wrap(err, responceUseCaseError)
	}
	response.ID = id
	img, err := r.mediaStore.Image().GetImage(response.UserImg)
	if err != nil {
		return nil, errors.Wrap(err, responceUseCaseError)
	}
	response.UserImg = string(img)
	return &response, nil
}

func (r *ResponseOrderUseCase) FindByVacancyID(id uint64) ([]model.ResponseOrder, error) {
	responses, err := r.store.ResponseOrder().FindByOrderId(id)
	if err != nil {
		return nil, errors.Wrap(err, responceUseCaseError)
	}
	for _, response := range responses {
		img, err := r.mediaStore.Image().GetImage(response.UserImg)
		if err != nil {
			return nil, errors.Wrap(err, responceUseCaseError)
		}
		response.UserImg = string(img)
	}
	if responses == nil {
		return []model.ResponseOrder{}, nil
	}
	return responses, nil
}

func (r *ResponseOrderUseCase) Change(response model.ResponseOrder) (*model.ResponseOrder, error) {
	changedResponse, err := r.store.ResponseOrder().Change(response)
	if err != nil {
		return nil, errors.Wrap(err, responceUseCaseError)
	}
	return changedResponse, nil
}

func (r *ResponseOrderUseCase) Delete(response model.ResponseOrder) error {
	err := r.store.ResponseOrder().Delete(response)
	if err != nil {
		return errors.Wrap(err, responceUseCaseError)
	}
	return nil
}
