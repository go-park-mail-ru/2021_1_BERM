
package teststore

import (
	"fl_ru/model"
	"fl_ru/store"
)

type Store struct {
	userRepository *UserRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[uint64]*model.User),
	}

	return s.userRepository
}
