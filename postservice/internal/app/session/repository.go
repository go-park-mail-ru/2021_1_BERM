package session

import (
	"golang.org/x/net/context"
	"post/internal/app/models"
)
//go:generate mockgen -destination=./mock/mock_repository.go -package=mock user/internal/app/session Repository
type Repository interface {
	Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error)
}
