package usecase

import (
	"context"
	"image/api"
	"image/internal/app/imgcreator"
	"image/internal/app/models"
)

type UseCase struct {
	UserRepo api.UserClient
}

func NewUseCase(userRepo api.UserClient) *UseCase {
	return &UseCase{
		UserRepo: userRepo,
	}
}

func (u *UseCase) SetImage(user models.UserImg) (models.UserImg, error) {
	imgUrl, err := imgcreator.CreateImg(user.Img)
	//TODO: обработка ошибки
	if err != nil {
		return models.UserImg{}, err
	}

	imgUrl, err = imgcreator.CropImg(imgUrl)
	//TODO: обработка ошибки
	if err != nil {
		return models.UserImg{}, err
	}
	_, err = u.UserRepo.SetImgUrl(context.Background(), &api.SetImgUrlRequest{Id: user.ID, ImgIrl: imgUrl})
	//TODO: обработка ошибки
	if err != nil {
		return models.UserImg{}, err
	}
	user.Img = imgUrl

	return user, nil
}
