package store

import "fl_ru/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(user *model.User) error
	Find(user *model.User) error
	ChangeUser(user *model.User) error
}

type SessionRepository interface {
	Create(session *model.Session) error
	Find(session *model.Session) error
}

type OrderRepository interface{
	Create(order*model.Order) error
	Find(order *model.Order) error
}