package impl

import (
	"context"
	"user/internal/app/models"
	"user/internal/session/repository"
)

type UseCase struct {
	sessionRepository repository.Repository
}

func (useCase *UseCase)Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error){
	u, err := useCase.sessionRepository.Check(sessionID)
	if err != nil{
		return nil, err
	}
	return u, err
}