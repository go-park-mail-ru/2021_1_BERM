package usecase

import (
	"context"
	"user/internal/app/models"
)

type UseCase interface {
	Create(user models.NewUser, ctx context.Context) (map[string]interface{}, error)
	Verification(email string, password string, ctx context.Context) (map[string]interface{}, error)
	GetById(ID uint64, ctx context.Context) (*models.UserInfo, error)
	Change(user models.ChangeUser, ctx context.Context) (map[string]interface{}, error)
	SetImg(ID uint64, img string, ctx context.Context) error
}
