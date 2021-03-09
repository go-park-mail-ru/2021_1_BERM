package tarantoolstore

import (
	"errors"
	"fl_ru/model"
	"github.com/tarantool/go-tarantool"
)

type SessionRepository struct {
	store *Store
}

func (s *SessionRepository)Create(session *model.Session) error{
	resp, err := s.store.conn.Insert("session", sessionToTarantoolData(session))
	println(resp)
	return err
}

func (s *SessionRepository)Find(session *model.Session) error{
	resp, err := s.store.conn.Select("session", "primary",
		0, 1,  tarantool.IterEq, []interface{}{
			session.SessionId,
		})
	if err != nil{
		return err
	}
	if len(resp.Tuples()) == 0{
		return errors.New("Not autorizate")
	}
	*session = *tarantoolDataToSession(resp.Tuples()[0])
	return nil
}

func sessionToTarantoolData(s *model.Session) []interface{}{
	return []interface{}{s.SessionId, s.UserId}
}

func tarantoolDataToSession(data []interface{}) *model.Session {
	s := &model.Session{}
	s.SessionId, _ = data[0].(string)
	s.UserId, _ = data[1].(uint64)
	return s
}