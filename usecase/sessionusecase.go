package usecase

import "FL_2/model"

type SessionUseCase interface {
	Create(u *model.User) (*model.Session, error)
	FindBySessionID(sessionID string)  (*model.Session, error)
}
