package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordSalt = "asdknj279312kasl0sshALkMnHG"
)

type User struct {
	Id          uint64   `json:"id,omitempty"`
	Email       string   `json:"email"`
	Password    string   `json:"password,omitempty"`
	UserName    string   `json:"user_name,omitempty"`
	FirstName   string   `json:"first_name,omitempty"`
	SecondName  string   `json:"second_name,omitempty"`
	Executor    bool     `json:"executor,omitempty"`
	Description string   `json:"description,omitempty,omitempty"`
	Specializes []string `json:"specializes,omitempty,omitempty"`
	ImgUrl      string   `json:"img_url,omitempty,omitempty"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 100)),
		validation.Field(&u.UserName, validation.Required),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.SecondName, validation.Required),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password, passwordSalt)
		if err != nil {
			return err
		}

		u.Password = enc
	}

	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password+passwordSalt)) == nil
}

func (u *User) Sanitize() {
	u.Password = ""
}
