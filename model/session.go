package model

const (
	cookieSalt = "wdsamlsdm2094dmfh"
)

type Session struct {
	SessionID string
	UserID    uint64
}

func (s *Session) BeforeChange() {
	s.SessionID, _ = EncryptString(s.SessionID, cookieSalt)
}
