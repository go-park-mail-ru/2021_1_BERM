package postgresstore

import (
	"FL_2/store"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db                 *sqlx.DB
	dsn                string
	userRepository     *UserRepository
	orderRepository    *OrderRepository
	vacancyRepository  *VacancyRepository
	responseRepository *ResponseRepository
}

func New(dsn string) *Store {
	return &Store{
		dsn: dsn,
	}
}

func (s *Store) Open() error {
	db, err := sqlx.Connect("postgres", s.dsn)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository == nil {
		s.orderRepository = &OrderRepository{
			store: s,
		}
	}

	return s.orderRepository
}

func (s *Store) Vacancy() store.VacancyRepository {
	if s.vacancyRepository == nil {
		s.vacancyRepository = &VacancyRepository{
			store: s,
		}
	}

	return s.vacancyRepository
}

func (s *Store) Response() store.ResponseRepository {
	if s.responseRepository == nil {
		s.responseRepository = &ResponseRepository{
			store: s,
		}
	}

	return s.responseRepository
}
