package implementation

import (
	"FL_2/model"
	"FL_2/store"
	errors "github.com/pkg/errors"
)

const (
	ErrInResponseVacancyUseCase = "Error in response vacancy use case"
)

type ResponseVacancyUseCase struct {
	store      store.Store
	mediaStore store.MediaStore
}

func (r *ResponseVacancyUseCase) Create(response model.ResponseVacancy) (*model.ResponseVacancy, error) {
	user, err := r.store.User().FindUserByID(response.UserID)
	if err != nil {
		return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	user.Specializes, err = r.store.User().FindSpecializesByUserID(response.UserID)
	if err != nil {
		return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	response.UserLogin = user.Login
	response.UserImg = user.Img
	id, err := r.store.ResponseVacancy().Create(response)
	response.ID = id
	if err != nil {
		return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	img, err := r.mediaStore.Image().GetImage(response.UserImg)
	if err != nil {
		return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	response.UserImg = string(img)
	return &response, nil
}

func (r *ResponseVacancyUseCase) FindByVacancyID(id uint64) ([]model.ResponseVacancy, error) {
	responses, err := r.store.ResponseVacancy().FindByVacancyID(id)
	if err != nil {
		return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	for _, response := range responses {
		img, err := r.mediaStore.Image().GetImage(response.UserImg)
		if err != nil {
			return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
		}
		response.UserImg = string(img)
	}
	if responses == nil {
		return []model.ResponseVacancy{}, nil
	}
	return responses, nil
}

func (r *ResponseVacancyUseCase) Change(response model.ResponseVacancy) (*model.ResponseVacancy, error) {
	changedResponse, err := r.store.ResponseVacancy().Change(response)
	if err != nil {
		return nil, errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	return changedResponse, nil
}

func (r *ResponseVacancyUseCase) Delete(response model.ResponseVacancy) error {
	err := r.store.ResponseVacancy().Delete(response)
	if err != nil {
		return errors.Wrap(err, ErrInResponseVacancyUseCase)
	}
	return nil
}
