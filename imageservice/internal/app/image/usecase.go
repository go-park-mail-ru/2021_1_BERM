package image

import "imageservice/internal/app/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock imageservice/internal/app/image UseCase
type UseCase interface {
	SetImage(user models.UserImg) (models.UserImg, error)
}
