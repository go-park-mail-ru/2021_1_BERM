package store

type Store interface {
	User() UserRepository
	Order() OrderRepository
	Vacancy() VacancyRepository
	Response() ResponseRepository
}
