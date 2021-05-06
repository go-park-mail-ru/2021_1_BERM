package image

import (
	"context"
	"imageservice/api"
	"imageservice/internal/app/imgcreator"
	"imageservice/internal/app/models"
)

type UseCase struct {
	UserRepo   api.UserClient
	ImgCreator imgcreator.ImgCreatorI
}

func NewUseCase(userRepo api.UserClient, imgCreator imgcreator.ImgCreatorI) *UseCase {
	return &UseCase{
		UserRepo:   userRepo,
		ImgCreator: imgCreator,
	}
}

func (u *UseCase) SetImage(user models.UserImg) (models.UserImg, error) {
	imgUrl, err := u.ImgCreator.CreateImg(user.Img)
	if err != nil {
		return models.UserImg{}, err
	}

	imgUrl, err = u.ImgCreator.CropImg(imgUrl)
	if err != nil {
		return models.UserImg{}, err
	}
	_, err = u.UserRepo.SetImgUrl(context.Background(), &api.SetImgUrlRequest{Id: user.ID, ImgIrl: imgUrl})
	if err != nil {
		return models.UserImg{}, err
	}
	user.Img = imgUrl

	return user, nil
}
