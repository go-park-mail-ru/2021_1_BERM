package usecase

import (
	"context"
)

type UseCase interface {
	Create(specialize string, ctx context.Context) (uint64, error);
	FindByUseID(ID uint64, ctx context.Context) ([]string, error)
	Remove(ID uint64, ctx context.Context) error
}
