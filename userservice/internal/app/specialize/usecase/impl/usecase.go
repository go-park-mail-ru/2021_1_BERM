package impl

import (
	"context"
	"user/internal/app/specialize/repository"
)

type UseCase struct {
	specializeRepository repository.Repository
}

func (useCase* UseCase)Create(specialize string, ctx context.Context) (uint64, error){
	ID, err := useCase.specializeRepository.Create(specialize, ctx)
	if err != nil{
		return 0, err
	}
	return ID, err
}

func (useCase* UseCase)Remove(ID uint64, ctx context.Context) error{
	err := useCase.Remove(ID, ctx)
	if err != nil{
		return  err
	}
	return err
}

func(useCase* UseCase)FindByUseID(ID uint64, ctx context.Context) ([]string, error){
	spec, err := useCase.specializeRepository.FindByUserID(ID, ctx)
	if err != nil{
		return nil, err
	}
	return spec, err
}