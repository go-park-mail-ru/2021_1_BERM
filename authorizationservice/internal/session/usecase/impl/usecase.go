package impl

import (
	"authorizationservice/internal/models"
	"authorizationservice/internal/session/repository"
	"authorizationservice/internal/tools"
	"context"
	"strconv"
	"time"
)

type UseCase struct {
	sessionRepository repository.Repository
}

func New(sessionRepository repository.Repository) *UseCase {
	return &UseCase{
		sessionRepository: sessionRepository,
	}
}

func (useCase *UseCase) Create(ID uint64, executor bool, ctx context.Context) (*models.Session, error) {
	session := &models.Session{
		SessionID: strconv.FormatUint(ID, 10) + time.Now().String(),
		UserId:    ID,
		Executor:  executor,
	}

	err := tools.BeforeCreate(session)
	if err != nil {
		return nil, err
	}
	err = useCase.sessionRepository.Store(session, ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (useCase *UseCase) Get(sessionID string, ctx context.Context) (*models.Session, error) {
	session, err := useCase.sessionRepository.Get(sessionID, ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (useCase *UseCase) Remove(sessionID string, ctx context.Context) error {
	err := useCase.sessionRepository.Remove(sessionID, ctx)
	if err != nil {
		return err
	}
	return nil
}
