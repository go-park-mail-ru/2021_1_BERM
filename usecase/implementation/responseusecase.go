package implementation

import (
	"FL_2/model"
	"FL_2/store"
)

type ResponseUseCase struct {
	store store.Store
	mediaStore store.MediaStore
}

func (r *ResponseUseCase)Create(response model.Response) (*model.Response, error){
	user, err := r.store.User().FindByID(response.UserID)
	if err != nil {
		return nil, err
	}
	response.UserLogin = user.Login
	response.UserImg = user.Img
	id, err := r.store.Response().Create(response)
	response.ID = id
	if err != nil{
		return nil, err
	}
	img, err :=r.mediaStore.Image().GetImage(response.UserImg)
	if err != nil{
		return nil, err
	}
	response.UserImg = string(img)
	return &response, nil
}

func (r *ResponseUseCase) FindByID(id uint64)  ([]model.Response, error){
	responses, err := r.store.Response().FindById(id)
	if err != nil{
		return nil, err
	}
	for _, response := range responses {
		img, err := r.mediaStore.Image().GetImage(response)
		if err != nil {
			return nil, err
		}
		response.UserImg = string(img)
	}
	return responses, nil
}
