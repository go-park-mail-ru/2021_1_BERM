package impl

import (
	"context"
	"user/internal/app/specialize/repository"
)

type UseCase struct {
	specializeRepository repository.Repository
}

func New(specializeRepository repository.Repository) *UseCase {
	return &UseCase{
		specializeRepository: specializeRepository,
	}
}

func (useCase *UseCase) Create(specialize string, ctx context.Context) (uint64, error) {
	ID, err := useCase.specializeRepository.Create(specialize, ctx)
	if err != nil {
		return 0, err
	}
	return ID, err
}

func (useCase *UseCase) Remove(ID uint64, ctx context.Context) error {
	err := useCase.Remove(ID, ctx)
	if err != nil {
		return err
	}
	return err
}

func (useCase *UseCase)AssociateWithUser(ID uint64, spec string, ctx context.Context)  error{
	specID, err := useCase.specializeRepository.FindByName(spec, ctx)
	if err != nil{
		specID, err = useCase.specializeRepository.Create(spec, ctx)
		if err != nil{
			return  err
		}
	}
	err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, ID, ctx)
	if err != nil{
		return  err
	}
	return nil
}
