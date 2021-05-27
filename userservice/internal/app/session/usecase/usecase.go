package usecase

import (
	"context"
	"user/internal/app/models"
	"user/internal/app/session"
)

type UseCase struct {
	SessionRepository session.Repository
}

func New(sessionRepository session.Repository) *UseCase {
	return &UseCase{
		SessionRepository: sessionRepository,
	}
}

func (useCase *UseCase) Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error) {
	u, err := useCase.SessionRepository.Check(sessionID, ctx)
	if err != nil {
		return nil, err
	}
	return u, err
}
