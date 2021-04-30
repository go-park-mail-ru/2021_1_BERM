package repository

import (
	"authorizationservice/internal/models"
	"context"
)

type Repository interface {
	Store(session *models.Session, ctx context.Context) error
	Get(sessionID string, ctx context.Context) (*models.Session, error)
	Remove(sessionID string, ctx context.Context) error
}
