package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/lib/pq"
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
	ID          uint64   `json:"id,omitempty" db:"id"`
	Email       string   `json:"email,omitempty" db:"email"`
	Password    string   `json:"password,omitempty" db:"password"`
	Login       string   `json:"login,omitempty" db:"login"`
	NameSurname string   `json:"name_surname,omitempty" db:"name_surname"`
	Executor    bool     `json:"executor,omitempty" db:"executor"`
	About       string   `json:"about,omitempty" db:"about"`
	Specializes pq.StringArray `json:"specializes,omitempty" db:"specializes"`
	Img         string   `json:"img,omitempty" db:"img"`
	Rating		int		 `json:"rating,omitempty" db:"rating"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(MinPswdLenght, MaxPswdLength)),
		validation.Field(&u.Login, validation.Required),
		validation.Field(&u.NameSurname, validation.Required),
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
