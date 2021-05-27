package image

import (
	"context"
	"imageservice/api"
	"imageservice/internal/app/imgcreator"
	"imageservice/internal/app/models"
)

type UseCase struct {
	UserRepo   api.UserClient
	ImgCreator imgcreator.CreatorI
}

func NewUseCase(userRepo api.UserClient, imgCreator imgcreator.CreatorI) *UseCase {
	return &UseCase{
		UserRepo:   userRepo,
		ImgCreator: imgCreator,
	}
}

func (u *UseCase) SetImage(user models.UserImg) (models.UserImg, error) {
	imgURL, err := u.ImgCreator.CreateImg(user.Img)
	if err != nil {
		return models.UserImg{}, err
	}

	imgURL, err = u.ImgCreator.CropImg(imgURL)
	if err != nil {
		return models.UserImg{}, err
	}
	_, err = u.UserRepo.SetImgUrl(context.Background(), &api.SetImgUrlRequest{Id: user.ID, ImgIrl: imgURL})
	if err != nil {
		return models.UserImg{}, err
	}
	user.Img = imgURL

	return user, nil
}
