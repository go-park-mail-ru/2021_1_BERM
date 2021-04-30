package impl

import (
	"authorizationservice/internal/models"
	"authorizationservice/internal/profile/repository"
	"context"
)

type UseCase struct {
	profileRepository repository.Repository
}

func New(profileRepository repository.Repository) *UseCase {
	return &UseCase{
		profileRepository: profileRepository,
	}
}

func (useCase *UseCase) Create(newUser models.NewUser, ctx context.Context) (*models.UserBasicInfo, error) {
	userBasicInfo, err := useCase.profileRepository.Create(newUser, ctx)
	if err != nil {
		return nil, err
	}
	return userBasicInfo, nil
}

func (useCase *UseCase) Authentication(email string, password string, ctx context.Context) (*models.UserBasicInfo, error) {
	userBasicInfo, err := useCase.profileRepository.Authentication(email, password, ctx)
	if err != nil {
		return nil, err
	}
	return userBasicInfo, nil
}
