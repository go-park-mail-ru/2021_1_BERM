package session

import (
	models2 "authorizationservice/internal/app/models"
	"context"
)

//go:generate mockgen -destination=./mock/mock_repsoitory.go -package=mock authorizationservice/internal/app/session Repository
type Repository interface {
	Store(session models2.Session, ctx context.Context) error
	Get(sessionID string, ctx context.Context) (*models2.Session, error)
	Remove(sessionID string, ctx context.Context) error
}
