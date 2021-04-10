package implementation

import (
	"FL_2/store"
	"FL_2/usecase"
)

type UseCase struct {
	orderUseCase 	*OrderUseCase
	userUseCase   	*UserUseCase
	mediaUseCase  	*MediaUseCase
	sessionUseCase  *SessionUseCase
	vacancyUseCase  *VacancyUseCase
	responseUseCase *ResponseUseCase
}

func New(store store.Store, cache store.Cash, mediaStore store.MediaStore) *UseCase{

	useCase := &UseCase{
  		orderUseCase: &OrderUseCase{
  			store: store,
		},
  		userUseCase: &UserUseCase{
  			store: store,
  			mediaStore: mediaStore,
		},
  		mediaUseCase: &MediaUseCase{
  			store: store,
  			mediaStore: mediaStore,
		},
  		sessionUseCase:	&SessionUseCase{
  			cache: cache,
		},
  		vacancyUseCase: &VacancyUseCase{
  			store: store,
		},
		responseUseCase: &ResponseUseCase{
  			store: store,
  			mediaStore: mediaStore,
		},
  	}

  return useCase
}

func (c *UseCase )Order() usecase.OrderUseCase{
	return c.orderUseCase
}

func (c *UseCase ) User()  usecase.UserUseCase{
	return c.userUseCase
}

func (c *UseCase ) Media() usecase.MediaUseCase{
	return c.mediaUseCase
}

func (c *UseCase )Session() usecase.SessionUseCase{
	return c.sessionUseCase
}

func (c *UseCase )Vacancy() usecase.VacancyUseCase{
	return c.vacancyUseCase
}
