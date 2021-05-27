package usecase_test

import (
	"authorizationservice/internal/app/models"
	"authorizationservice/internal/app/profile/mock"
	profUCase "authorizationservice/internal/app/profile/usecase"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

//Проверка созданияю отзыва
func TestCreateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	newUser := &models.NewUser{}
	basicUseInfi := &models.UserBasicInfo{}
	mockProfileRep := mock.NewMockRepository(ctrl)
	mockProfileRep.EXPECT().Create(*newUser, ctx).Times(1).Return(basicUseInfi, nil)
	useCase := profUCase.UseCase{
		ProfileRepository: mockProfileRep,
	}

	_, err := useCase.Create(*newUser, ctx)
	require.NoError(t, err)
}

func TestCreateProfileErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	newUser := &models.NewUser{}
	basicUseInfi := &models.UserBasicInfo{}
	mockProfileRep := mock.NewMockRepository(ctrl)
	mockProfileRep.EXPECT().Create(*newUser, ctx).Times(1).Return(basicUseInfi, errors.New("err"))
	useCase := profUCase.UseCase{
		ProfileRepository: mockProfileRep,
	}

	_, err := useCase.Create(*newUser, ctx)
	require.Error(t, err)
}

func TestAuthenticationProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	basicUseInfi := &models.UserBasicInfo{}
	mockProfileRep := mock.NewMockRepository(ctrl)
	mockProfileRep.EXPECT().Authentication("1", "1", ctx).Times(1).Return(basicUseInfi, nil)
	useCase := profUCase.UseCase{
		ProfileRepository: mockProfileRep,
	}

	_, err := useCase.Authentication("1", "1", ctx)
	require.NoError(t, err)
}

func TestAuthenticationProfileErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	basicUseInfi := &models.UserBasicInfo{}
	mockProfileRep := mock.NewMockRepository(ctrl)
	mockProfileRep.EXPECT().Authentication("1", "1", ctx).Times(1).Return(basicUseInfi, errors.New("Err"))
	useCase := profUCase.UseCase{
		ProfileRepository: mockProfileRep,
	}

	_, err := useCase.Authentication("1", "1", ctx)
	require.Error(t, err)
}

func TestNewProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := mock.NewMockRepository(ctrl)

	uc := profUCase.New(mockProfileRepo)

	require.Equal(t, uc.ProfileRepository, mockProfileRepo)
}
