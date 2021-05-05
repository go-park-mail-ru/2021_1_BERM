package repository

import (
	"authorizationservice/internal/models"
	"context"
)

type Repository interface {
	Create(newUser models.NewUser, ctx context.Context) (*models.UserBasicInfo, error)
	Authentication(email string, password string, ctx context.Context) (*models.UserBasicInfo, error)
}
