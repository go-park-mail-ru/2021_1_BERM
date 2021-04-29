package repository

import (
	"context"
	"github.com/lib/pq"
)

type Repository interface {
	Create(specialize string, ctx context.Context) (uint64, error)
	FindByUserID(userID uint64, ctx context.Context) (pq.StringArray, error)
	AssociateSpecializationWithUser(specId uint64, userId uint64, ctx context.Context)  error
	FindById(ID uint64, ctx context.Context) (string, error)
	FindByName(spec string,  ctx context.Context) (uint64, error)
	Remove(ID uint64, ctx context.Context) error

}
