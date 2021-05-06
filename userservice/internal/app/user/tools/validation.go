package tools

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"user/internal/app/models"
)

const (
	minPasswordLength = 5
	maxPasswordLength = 30
)

func ValidationCreateUser(user *models.NewUser) error {
	err := validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.NameSurname, validation.Required),
	)
	if err != nil {
		return err
	}
	return nil
}

func ValidationChangeUser(user models.ChangeUser) error {

	err := validation.ValidateStruct(
		&user,
		validation.Field(&user.Email, is.Email),
		validation.Field(&user.Password, validation.Length(minPasswordLength, maxPasswordLength)),
	)
	if err != nil {
		return err
	}
	return nil
}