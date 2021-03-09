package model

import "golang.org/x/crypto/bcrypt"

const(
	cookieSalt = "wdsamlsdm2094dmfh"
)
type Session struct {
	SessionId string
	UserId uint64
}


func (s *Session) BeforeChange(){
	s.SessionId, _ = encryptString(s.SessionId, cookieSalt)
}

func (u *Session) CompareSessionId(sessionId string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.SessionId), []byte(sessionId)) == nil
}
