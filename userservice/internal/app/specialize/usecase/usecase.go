package usecase

import (
	"context"
)

type UseCase interface {
	Create(specialize string, ctx context.Context) (uint64, error)
	AssociateWithUser(ID uint64, spec string, ctx context.Context)  error
	Remove(ID uint64, spec string, ctx context.Context) error
}
