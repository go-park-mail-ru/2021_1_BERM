package usecase

import "FL_2/model"

type MediaUseCase interface {
	GetImage(imageInfo interface{}) (*model.User, error)
	SetImage(imageInfo interface{}, image []byte) (*model.User, error)
}
