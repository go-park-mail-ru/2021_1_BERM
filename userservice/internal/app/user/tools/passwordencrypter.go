package tools

import "user/internal/app/models"

//go:generate mockgen -destination=../mock/mock_tools.go -package=mock user/internal/app/user/tools PasswordEncrypter
type PasswordEncrypter interface {
	CompPass(passHash []byte, plainPassword string) bool
	BeforeCreate(user models.NewUser) (models.NewUser, error)
	BeforeChange(user models.ChangeUser) (models.ChangeUser, error)
}
