package session

import (
	"context"
	"post/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock post/internal/app/session Repository
type Repository interface {
	Check(ctx context.Context, sessionID string) (*models.UserBasicInfo, error)
}
