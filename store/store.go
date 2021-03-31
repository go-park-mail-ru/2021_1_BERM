package store

type Store interface {
	User() UserRepository
	Order() OrderRepository
}
