package usecase

import (
	"authorizationservice/internal/app/models"
	"authorizationservice/internal/app/profile/mock"
	"context"
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
	useCase := UseCase{
		profileRepository: mockProfileRep,
	}

	_, err := useCase.Create(*newUser, ctx)
	require.NoError(t, err)
}


func TestAuthenticationProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	basicUseInfi := &models.UserBasicInfo{}
	mockProfileRep := mock.NewMockRepository(ctrl)
	mockProfileRep.EXPECT().Authentication("1", "1", ctx).Times(1).Return(basicUseInfi, nil)
	useCase := UseCase{
		profileRepository: mockProfileRep,
	}

	_, err := useCase.Authentication("1", "1", ctx)
	require.NoError(t, err)
}


func TestNewProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := mock.NewMockRepository(ctrl)

	uc := New(mockProfileRepo)

	require.Equal(t, uc.profileRepository, mockProfileRepo)
}