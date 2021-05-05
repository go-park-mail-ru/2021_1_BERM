package usecase

import (
	"authorizationservice/internal/models"
	"context"
)

type UseCase interface {
	Create(newUser models.NewUser, ctx context.Context) (*models.UserBasicInfo, error)
	Authentication(email string, password string, ctx context.Context) (*models.UserBasicInfo, error)
}
