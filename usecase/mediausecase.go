package usecase

import "FL_2/model"

type MediaUseCase interface {
	SetImage(imageInfo interface{}, image []byte) (*model.User, error)
}
