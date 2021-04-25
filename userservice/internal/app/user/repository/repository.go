package repository

import "user/internal/app/models"

type Repository interface {
	Create(user *models.NewUser) (uint64, error)
	Change(user *models.ChangeUser) error
	FindUserByID(ID uint64) (*models.UserInfo, error)
	FindUserByEmail(email string) (*models.UserInfo, error)
}

