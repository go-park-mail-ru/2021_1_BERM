package store

import (
	"FL_2/model"
	"github.com/lib/pq"
)

type UserRepository interface {
	AddUserSpec(userID uint64, specID uint64) error
	FindSpecializeByName(specName string) (model.Specialize, error)
	IsUserHaveSpec(specID uint64, userID uint64) (bool, error)
	AddSpec(specName string) (uint64, error)
	AddUser(user model.User) (uint64, error)
	FindUserByEmail(email string) (*model.User, error)
	FindSpecializesByUserEmail(email string) (pq.StringArray, error)
	FindUserByID(id uint64) (*model.User, error)
	FindSpecializesByUserID(id uint64) (pq.StringArray, error)
	ChangeUser(user model.User) (*model.User, error)
	DelSpecialize(specID uint64, userID uint64) error
}
