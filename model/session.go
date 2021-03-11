package model

const (
	cookieSalt = "wdsamlsdm2094dmfh"
)

type Session struct {
	SessionId string
	UserId    uint64
}

func (s *Session) BeforeChange() {
	s.SessionId, _ = encryptString(s.SessionId, cookieSalt)
}
