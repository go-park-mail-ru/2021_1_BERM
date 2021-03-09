package store

type Store interface {
	User() UserRepository
	Session() SessionRepository
	Order() OrderRepository
}
