package store

import "FL_2/model"

type UserRepository interface {
	Create(user *model.User) (uint64, error)
	FindByEmail(email string) (*model.User, error)
	FindById(id uint64) (*model.User, error)
	ChangeUser(user *model.User) (*model.User, error)
}
