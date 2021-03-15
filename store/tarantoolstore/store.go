package tarantoolstore

import (
	"fl_ru/store"
	"github.com/tarantool/go-tarantool"
)

type Store struct {
	conn              *tarantool.Connection
	UserRepository    *UserRepository
	SessionRepository *SessionRepository
	OrderRepository   *OrderRepository
}

func New(dbURL string) (*Store, error) {
	conn, err := newTarantoolConnect(dbURL)
	if err != nil {
		return nil, err
	}

	return &Store{
		conn: conn,
	}, nil
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}
	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}

func (s *Store) Session() store.SessionRepository {
	if s.SessionRepository != nil {
		return s.SessionRepository
	}
	s.SessionRepository = &SessionRepository{
		store: s,
	}

	return s.SessionRepository
}

func (s *Store) Order() store.OrderRepository {
	if s.OrderRepository != nil {
		return s.OrderRepository
	}
	s.OrderRepository = &OrderRepository{
		store: s,
	}

	return s.OrderRepository
}

func newTarantoolConnect(dbURL string) (*tarantool.Connection, error) {
	opts := tarantool.Opts{User: "guest"}
	db, err := tarantool.Connect(dbURL, opts)

	return db, err
}
