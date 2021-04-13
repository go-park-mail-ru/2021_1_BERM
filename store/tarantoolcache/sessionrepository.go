package tarantoolcache

import (
	"FL_2/model"
	"github.com/pkg/errors"
	"github.com/tarantool/go-tarantool"
)

type SessionRepository struct {
	cache *Cache
}

func (s *SessionRepository) Create(session *model.Session) error {
	_, err := s.cache.conn.Insert("session", sessionToTarantoolData(session))
	if err != nil {
		return errors.Wrap(err, sessionSourceError)
	}
	return err
}

func (s *SessionRepository) Find(session *model.Session) error {
	resp, err := s.cache.conn.Select("session", "primary",
		0, 1, tarantool.IterEq, []interface{}{
			session.SessionID,
		})
	if err != nil {
		return errors.Wrap(err, sessionSourceError)
	}
	if len(resp.Tuples()) == 0 {
		return errors.Wrap(NotAuthorized, sessionSourceError)
	}
	*session = *tarantoolDataToSession(resp.Tuples()[0])
	return nil
}

func sessionToTarantoolData(s *model.Session) []interface{} {
	return []interface{}{s.SessionID, s.UserID}
}

func tarantoolDataToSession(data []interface{}) *model.Session {
	s := &model.Session{}
	s.SessionID, _ = data[0].(string)
	s.UserID, _ = data[1].(uint64)
	return s
}
