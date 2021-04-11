package usecase

import "FL_2/model"

type UserUseCase interface {
	Create(user *model.User) error
	UserVerification(email string, password string) (*model.User, error)
	FindByID(id uint64) (*model.User, error)
	ChangeUser(user model.User) (*model.User, error)
	AddSpecialize(specName string, userID uint64) error
	DelSpecialize(specName string, userID uint64) error
}
