package usecase

import (
	"context"
	"post/internal/app/models"
	"post/internal/app/session"
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
	u, err := useCase.SessionRepository.Check(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return u, err
}
