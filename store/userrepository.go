package store

import "FL_2/model"

type UserRepository interface {
	Create(user model.User) (uint64, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint64) (*model.User, error)
	ChangeUser(user model.User) (*model.User, error)
	AddSpecialize(specName string, userID uint64) error
	DelSpecialize(specName string, userID uint64) error
}
