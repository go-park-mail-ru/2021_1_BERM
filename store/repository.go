package store

import (
	"fl_ru/model"
)

type UserRepository interface {
	Create(user *model.User) (uint64, error)
	FindByEmail(email string) (*model.User, error)
	FindById(id uint64) (*model.User, error)
	ChangeUser(user *model.User) (*model.User, error)
}

type OrderRepository interface {
	Create(order *model.Order) error
	Find(order *model.Order) error
}
