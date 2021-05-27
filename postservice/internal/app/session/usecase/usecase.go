package usecase

import (
	"context"
	"post/internal/app/models"
	"post/internal/app/session"
)

type UseCase struct {
	sessionRepository session.Repository
}

func New(sessionRepository session.Repository) *UseCase {
	return &UseCase{
		sessionRepository: sessionRepository,
	}
}

func (useCase *UseCase) Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error) {
	u, err := useCase.sessionRepository.Check(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return u, err
}
