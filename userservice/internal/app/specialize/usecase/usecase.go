package usecase

import (
	"context"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"user/internal/app/specialize"
)

type UseCase struct {
	SpecializeRepository specialize.Repository
}

func New(specializeRepository specialize.Repository) *UseCase {
	return &UseCase{
		SpecializeRepository: specializeRepository,
	}
}

func (useCase *UseCase) Create(specialize string, ctx context.Context) (uint64, error) {
	ID, err := useCase.SpecializeRepository.Create(specialize, ctx)
	if err != nil {
		return 0, err
	}
	return ID, err
}

func (useCase *UseCase) Remove(ID uint64, spec string, ctx context.Context) error {
	specID, err := useCase.SpecializeRepository.FindByName(spec, ctx)
	if err != nil {
		return err
	}
	err = useCase.SpecializeRepository.RemoveAssociateSpecializationWithUser(specID, ID, ctx)
	if err != nil {
		return err
	}
	return err
}

func (useCase *UseCase) AssociateWithUser(ID uint64, spec string, ctx context.Context) error {
	specID, err := useCase.SpecializeRepository.FindByName(spec, ctx)
	if err != nil {
		specID, err = useCase.SpecializeRepository.Create(spec, ctx)
		if err != nil {
			return err
		}
	}
	err = useCase.SpecializeRepository.AssociateSpecializationWithUser(specID, ID, ctx)
	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			return errors.New("Duplicate spec")
		}
	}
	if err != nil {
		return err
	}
	return nil
}
