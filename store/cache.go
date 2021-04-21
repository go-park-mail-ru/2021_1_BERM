package store

//go:generate mockgen -destination=mock/mock_cache_repo.go -package=mock FL_2/store Caсhe
type Caсhe interface {
	Session() SessionRepository
}
