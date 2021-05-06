package session

import (
	"context"
	"user/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock user/internal/app/session UseCase
type UseCase interface {
	Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error)
}
