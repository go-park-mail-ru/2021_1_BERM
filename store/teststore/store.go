package teststore

import (
	"fl_ru/model"
	"fl_ru/store"
)

type Store struct {
	UserRepository    *UserRepository
	SessionRepository *SessionRepository
	OrderRepository   *OrderRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
		users: make(map[uint64]model.User),
	}

	return s.UserRepository
}

func (s *Store) Session() store.SessionRepository {
	if s.SessionRepository != nil {
		return s.SessionRepository
	}

	s.SessionRepository = &SessionRepository{
		store:    s,
		sessions: make(map[string]*model.Session),
	}

	return s.SessionRepository
}

func (s *Store) Order() store.OrderRepository {
	if s.OrderRepository != nil {
		return s.OrderRepository
	}

	s.OrderRepository = &OrderRepository{
		store: s,
		order: make(map[uint64]*model.Order),
	}

	return s.OrderRepository
}
