package store

//go:generate mockgen -destination=mock/mock_store_store.go -package=mock FL_2/store Store
type Store interface {
	User() UserRepository
	Order() OrderRepository
	Vacancy() VacancyRepository
	ResponseOrder() ResponseOrderRepository
	ResponseVacancy() ResponseVacancyRepository
}
