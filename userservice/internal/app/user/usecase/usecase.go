package usecase

import "user/internal/app/models"

type UseCase interface {
	Create(user models.NewUser) (map[string]interface{}, error)
	Verification(email string, password string) (map[string]interface{}, error)
	GetById(ID uint64) (*models.UserInfo, error)
	Change(user models.ChangeUser) (map[string]interface{}, error)
}
