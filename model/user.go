package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordSalt = "asdknj279312kasl0sshALkMnHG"
)

const (
	MinPswdLenght int = 5
	MaxPswdLength int = 300
)

type User struct {
	ID          uint64   `json:"id,omitempty"`
	Email       string   `json:"email"`
	Password    string   `json:"password,omitempty"`
	UserName    string   `json:"user_name,omitempty"`
	FirstName   string   `json:"first_name,omitempty"`
	SecondName  string   `json:"second_name,omitempty"`
	Executor    bool     `json:"executor,omitempty"`
	About       string   `json:"about,omitempty"`
	Specializes []string `json:"specializes,omitempty"`
	ImgURL      string   `json:"img_url,omitempty"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(MinPswdLenght, MaxPswdLength)),
		validation.Field(&u.UserName, validation.Required),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.SecondName, validation.Required),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := EncryptString(u.Password, passwordSalt)
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
