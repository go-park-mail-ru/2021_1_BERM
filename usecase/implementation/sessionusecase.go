package implementation

import (
	"FL_2/model"
	"FL_2/store"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	cookieSalt          = "wdsamlsdm2094dmfh"
	sessionUseCaseError = "Session use case error"
)

type SessionUseCase struct {
	cache store.Caсhe
}

func (s *SessionUseCase) Create(u *model.User) (*model.Session, error) {
	session := &model.Session{
		SessionID: u.Email + time.Now().String(),
		UserID:    u.ID,
	}

	err := s.beforeCreate(session)
	if err != nil {
		return nil, errors.Wrap(err, sessionUseCaseError)
	}
	if err = s.cache.Session().Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionUseCase) FindBySessionID(sessionID string) (*model.Session, error) {
	session := &model.Session{
		SessionID: sessionID,
	}
	err := s.cache.Session().Find(session)
	if err != nil {
		return nil, errors.Wrap(err, sessionUseCaseError)
	}
	session.SessionID = ""
	return session, err
}

func (s *SessionUseCase) encryptString(password string, salt string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(err, sessionUseCaseError)
	}
	return string(b), nil
}

func (s *SessionUseCase) beforeCreate(session *model.Session) error {
	var err error
	session.SessionID, err = s.encryptString(session.SessionID, cookieSalt)
	return err
}
