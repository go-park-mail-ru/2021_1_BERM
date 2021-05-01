package image

import "image/internal/app/models"

type UseCase interface {
	SetImage(user models.UserImg) (models.UserImg, error)
}
