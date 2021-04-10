package usecase

type UseCase interface {
	Order() OrderUseCase
	User() 	UserUseCase
	Media() MediaUseCase
	Session() SessionUseCase
	Vacancy() VacancyUseCase
	Response() ResponseUseCase
}
