package implementation

import (
	"FL_2/model"
	"FL_2/store"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const(
	cookieSalt = "wdsamlsdm2094dmfh"
)

type SessionUseCase struct {
	cache store.Cash
}


func (s *SessionUseCase)Create(u *model.User) (*model.Session, error){
	session := &model.Session{
		SessionId: u.Email + time.Now().String(),
		UserId: u.ID,
	}

	err := s.beforeCreate(session)
	if err != nil {
		return nil, err
	}
	if err = s.cache.Session().Create(session); err != nil {
		return nil, err
	}
	return session, nil;
}

func (s *SessionUseCase)FindBySessionID(sessionID string)  (*model.Session, error) {
	session := &model.Session{}
	err := s.cache.Session().Find(session)
	if err != nil{
		return nil, err
	}
	return session, err
}

func (s *SessionUseCase) encryptString(password string, salt string) (string, error){
	b, err := bcrypt.GenerateFromPassword([]byte(password + salt), bcrypt.MinCost)
	if err != nil{
		return "", err
	}
	return string(b), nil
}

func (s *SessionUseCase) beforeCreate(session *model.Session) error{
	var err error
	session.SessionId, err = s.encryptString(session.SessionId, cookieSalt)
	return err
}
//
//func (s *SessionUseCase) compareSessionId(session model.Session, sessionId string) bool {
//	return bcrypt.CompareHashAndPassword([]byte(session.SessionId), []byte(sessionId)) == nil
//}
