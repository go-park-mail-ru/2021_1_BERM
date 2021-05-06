package specialize

import (
	"context"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock user/internal/app/specialize UseCase
type UseCase interface {
	Create(specialize string, ctx context.Context) (uint64, error)
	AssociateWithUser(ID uint64, spec string, ctx context.Context) error
	Remove(ID uint64, spec string, ctx context.Context) error
}
