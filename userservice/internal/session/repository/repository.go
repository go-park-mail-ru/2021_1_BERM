package repository

import (
	"golang.org/x/net/context"
	"user/internal/app/models"
)

type Repository interface {
	Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error)
}