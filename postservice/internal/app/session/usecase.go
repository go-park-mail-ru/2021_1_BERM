package session

import (
	"context"
	"post/internal/app/models"
)

type UseCase interface {
	Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error)
}
