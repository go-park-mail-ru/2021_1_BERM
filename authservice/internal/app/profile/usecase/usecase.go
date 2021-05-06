package usecase

import (
	models2 "authorizationservice/internal/app/models"
	profile2 "authorizationservice/internal/app/profile"
	"context"
)

type UseCase struct {
	profileRepository profile2.Repository
}

func New(profileRepository profile2.Repository) *UseCase {
	return &UseCase{
		profileRepository: profileRepository,
	}
}

func (useCase *UseCase) Create(newUser models2.NewUser, ctx context.Context) (*models2.UserBasicInfo, error) {
	userBasicInfo, err := useCase.profileRepository.Create(newUser, ctx)
	if err != nil {
		return nil, err
	}
	return userBasicInfo, nil
}

func (useCase *UseCase) Authentication(email string, password string, ctx context.Context) (*models2.UserBasicInfo, error) {
	userBasicInfo, err := useCase.profileRepository.Authentication(email, password, ctx)
	if err != nil {
		return nil, err
	}
	return userBasicInfo, nil
}
