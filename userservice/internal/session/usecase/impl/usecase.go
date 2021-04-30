package impl

import (
	"context"
	"user/internal/app/models"
	"user/internal/session/repository"
)

type UseCase struct {
	sessionRepository repository.Repository
}

func New(sessionRepository repository.Repository) *UseCase{
	return &UseCase{
		sessionRepository: sessionRepository,
	}
}

func (useCase *UseCase)Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error){
	u, err := useCase.sessionRepository.Check(sessionID, ctx)
	if err != nil{
		return nil, err
	}
	return u, err
}