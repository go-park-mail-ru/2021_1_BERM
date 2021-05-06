package usecase

import (
	"context"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"user/internal/app/specialize"
)

type UseCase struct {
	specializeRepository specialize.Repository
}

func New(specializeRepository specialize.Repository) *UseCase {
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

func (useCase *UseCase) Remove(ID uint64, spec string, ctx context.Context) error {
	specID, err := useCase.specializeRepository.FindByName(spec, ctx)
	if err != nil {
		return err
	}
	err = useCase.specializeRepository.RemoveAssociateSpecializationWithUser(specID, ID, ctx)
	if err != nil {
		return err
	}
	return err
}

func (useCase *UseCase) AssociateWithUser(ID uint64, spec string, ctx context.Context) error {
	specID, err := useCase.specializeRepository.FindByName(spec, ctx)
	if err != nil {
		specID, err = useCase.specializeRepository.Create(spec, ctx)
		if err != nil {
			return err
		}
	}
	err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, ID, ctx)
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
