package image_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"imageservice/api"
	"imageservice/internal/app/image/mock"
	uCase "imageservice/internal/app/image/usecase"
	mockCreator "imageservice/internal/app/imgcreator/mock"
	"imageservice/internal/app/models"
	"testing"
)

func TestUseCase_SetImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCr := mockCreator.NewMockImgCreatorI(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := uCase.NewUseCase(mockUserRepo, mockCr)

	imgUrl := "kek.jpeg"
	user := models.UserImg{
		Img: "/9j/4AAQSkZJRgA",
		ID:  1,
	}
	expectUser := models.UserImg{
		Img: imgUrl,
		ID:  1,
	}
	mockUserRepo.EXPECT().
		SetImgUrl(context.Background(), &api.SetImgUrlRequest{Id: user.ID, ImgIrl: imgUrl}).
		Times(1).
		Return(&api.UserInfoResponse{}, nil)
	mockCr.EXPECT().
		CreateImg(user.Img).
		Times(1).
		Return(imgUrl, nil)
	mockCr.EXPECT().
		CropImg(imgUrl).
		Times(1).
		Return(imgUrl, nil)
	resp, err := useCase.SetImage(user)

	require.NoError(t, err)
	require.Equal(t, expectUser.ID, resp.ID)
}

func TestUseCase_SetImageErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCr := mockCreator.NewMockImgCreatorI(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := uCase.NewUseCase(mockUserRepo, mockCr)

	imgUrl := "kek.jpeg"
	user := models.UserImg{
		Img: "/9j/4AAQSkZJRgA",
		ID:  1,
	}

	mockCr.EXPECT().
		CreateImg(user.Img).
		Times(1).
		Return(imgUrl, errors.New("Err"))

	_, err := useCase.SetImage(user)

	require.Error(t, err)
}

func TestUseCase_SetImageErr2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCr := mockCreator.NewMockImgCreatorI(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := uCase.NewUseCase(mockUserRepo, mockCr)

	imgUrl := "kek.jpeg"
	user := models.UserImg{
		Img: "/9j/4AAQSkZJRgA",
		ID:  1,
	}

	mockCr.EXPECT().
		CreateImg(user.Img).
		Times(1).
		Return(imgUrl, nil)
	mockCr.EXPECT().
		CropImg(imgUrl).
		Times(1).
		Return(imgUrl, errors.New("err"))
	_, err := useCase.SetImage(user)

	require.Error(t, err)
}

func TestUseCase_SetImage3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCr := mockCreator.NewMockImgCreatorI(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := uCase.NewUseCase(mockUserRepo, mockCr)

	imgUrl := "kek.jpeg"
	user := models.UserImg{
		Img: "/9j/4AAQSkZJRgA",
		ID:  1,
	}
	mockUserRepo.EXPECT().
		SetImgUrl(context.Background(), &api.SetImgUrlRequest{Id: user.ID, ImgIrl: imgUrl}).
		Times(1).
		Return(&api.UserInfoResponse{}, errors.New("Err"))
	mockCr.EXPECT().
		CreateImg(user.Img).
		Times(1).
		Return(imgUrl, nil)
	mockCr.EXPECT().
		CropImg(imgUrl).
		Times(1).
		Return(imgUrl, nil)
	_, err := useCase.SetImage(user)

	require.Error(t, err)
}
