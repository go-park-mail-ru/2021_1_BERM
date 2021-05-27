package profile

import (
	models2 "authorizationservice/internal/app/models"
	"context"
)

//nolint:lll    //go:generate mockgen -destination=./mock/mock_usecase.go -package=mock authorizationservice/internal/app/profile UseCase
type UseCase interface {
	Create(newUser models2.NewUser, ctx context.Context) (*models2.UserBasicInfo, error)
	Authentication(email string, password string, ctx context.Context) (*models2.UserBasicInfo, error)
}
