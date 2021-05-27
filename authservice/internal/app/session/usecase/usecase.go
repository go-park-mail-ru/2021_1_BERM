package usecase

import (
	models2 "authorizationservice/internal/app/models"
	session2 "authorizationservice/internal/app/session"
	tools2 "authorizationservice/internal/app/session/tools"
	sessiontools2 "authorizationservice/internal/app/session/tools/sessiontools"
	"context"
)

type UseCase struct {
	SessionRepository session2.Repository
	Tools             tools2.SessionTools
}

func New(sessionRepository session2.Repository) *UseCase {
	return &UseCase{
		SessionRepository: sessionRepository,
		Tools:             &sessiontools2.SessionTools{},
	}
}

func (useCase *UseCase) Create(ID uint64, executor bool, ctx context.Context) (*models2.Session, error) {
	session := models2.Session{
		UserId:   ID,
		Executor: executor,
	}
	var err error
	session, err = useCase.Tools.BeforeCreate(session)
	if err != nil {
		return nil, err
	}
	err = useCase.SessionRepository.Store(session, ctx)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (useCase *UseCase) Get(sessionID string, ctx context.Context) (*models2.Session, error) {
	session, err := useCase.SessionRepository.Get(sessionID, ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (useCase *UseCase) Remove(sessionID string, ctx context.Context) error {
	err := useCase.SessionRepository.Remove(sessionID, ctx)
	if err != nil {
		return err
	}
	return nil
}
