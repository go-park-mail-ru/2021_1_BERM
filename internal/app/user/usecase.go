package user

import "ff/internal/app/models"

type UserUseCase interface {
	Create(user *models.User) error
	UserVerification(email string, password string) (*models.User, error)
	FindByID(id uint64) (*models.User, error)
	ChangeUser(user models.User) (*models.User, error)
	AddSpecialize(specName string, userID uint64) error
	DelSpecialize(specName string, userID uint64) error
}

