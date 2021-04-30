package repository

import (
	"context"
	"user/internal/app/models"
)

type Repository interface {
	Create(user *models.NewUser, ctx context.Context) (uint64, error)
	Change(user *models.ChangeUser, ctx context.Context) error
	FindUserByID(ID uint64, ctx context.Context) (*models.UserInfo, error)
	FindUserByEmail(email string, ctx context.Context) (*models.UserInfo, error)
	SetUserImg(ID uint64, img string, ctx context.Context) error
}

