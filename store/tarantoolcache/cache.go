package tarantoolcache

import (
	"FL_2/store"
	"github.com/tarantool/go-tarantool"
)

type Cache struct {
	conn              *tarantool.Connection
	SessionRepository *SessionRepository
}

func New(dbURL string) (*Cache, error) {
	conn, err := newTarantoolConnect(dbURL)
	if err != nil {
		return nil, err
	}
	return &Cache{
		conn: conn,
	}, nil
}

func (s *Cache) Session() store.SessionRepository {
	if s.SessionRepository != nil {
		return s.SessionRepository
	}
	s.SessionRepository = &SessionRepository{
		cache: s,
	}
	return s.SessionRepository
}

func newTarantoolConnect(dbURL string) (*tarantool.Connection, error) {
	opts := tarantool.Opts{User: "guest"}
	db, err := tarantool.Connect(dbURL, opts)
	return db, err
}
