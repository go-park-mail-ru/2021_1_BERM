package profile

import (
	models2 "authorizationservice/internal/app/models"
	"context"
)

//go:generate mockgen -destination=./mock/mock_repsoitory.go -package=mock authorizationservice/internal/app/profile Repository
type Repository interface {
	Create(newUser models2.NewUser, ctx context.Context) (*models2.UserBasicInfo, error)
	Authentication(email string, password string, ctx context.Context) (*models2.UserBasicInfo, error)
}
