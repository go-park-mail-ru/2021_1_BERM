package usecase

import (
	"authorizationservice/internal/models"
	"context"
)

type UseCase interface {
	Create(ID uint64, executor bool, ctx context.Context) (*models.Session, error)
	Get(sessionID string, ctx context.Context) (*models.Session, error)
	Remove(sessionID string, ctx context.Context) error
}
