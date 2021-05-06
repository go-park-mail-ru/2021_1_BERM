package session

import (
	models2 "authorizationservice/internal/app/models"
	"context"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock authorizationservice/internal/app/session UseCase
type UseCase interface {
	Create(ID uint64, executor bool, ctx context.Context) (*models2.Session, error)
	Get(sessionID string, ctx context.Context) (*models2.Session, error)
	Remove(sessionID string, ctx context.Context) error
}
