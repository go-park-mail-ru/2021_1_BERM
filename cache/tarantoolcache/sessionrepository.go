package tarantoolcache

import (
	"FL_2/model"
	"errors"
	"github.com/tarantool/go-tarantool"
)

type SessionRepository struct {
	cache *Cache
}

func (s *SessionRepository) Create(session *model.Session) error {
	resp, err := s.cache.conn.Insert("session", sessionToTarantoolData(session))
	println(resp)
	return err
}

func (s *SessionRepository) Find(session *model.Session) error {
	resp, err := s.cache.conn.Select("session", "primary",
		0, 1, tarantool.IterEq, []interface{}{
			session.SessionId,
		})
	if err != nil {
		return err
	}
	if len(resp.Tuples()) == 0 {
		return errors.New("Not autorizate")
	}
	*session = *tarantoolDataToSession(resp.Tuples()[0])
	return nil
}

func sessionToTarantoolData(s *model.Session) []interface{} {
	return []interface{}{s.SessionId, s.UserId}
}

func tarantoolDataToSession(data []interface{}) *model.Session {
	s := &model.Session{}
	s.SessionId, _ = data[0].(string)
	s.UserId, _ = data[1].(uint64)
	return s
}