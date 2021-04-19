package models

import (
	"github.com/lib/pq"
)

type User struct {
	ID              uint64         `json:"id,omitempty" db:"id"`
	Email           string         `json:"email,omitempty" db:"email"`
	Password        string         `json:"password,omitempty" db:"-"`
	NewPassword     string         `json:"new_password,omitempty" db:"-"`
	EncryptPassword []byte         `json:"-" db:"password"`
	Login           string         `json:"login,omitempty" db:"login"`
	NameSurname     string         `json:"name_surname,omitempty" db:"name_surname"`
	Executor        bool           `json:"executor,omitempty" db:"executor"`
	About           string         `json:"about,omitempty" db:"about"`
	Specializes     pq.StringArray `json:"specializes,omitempty" db:"specializes,omitempty"`
	Img             string         `json:"img,omitempty" db:"img"`
	Rating          int            `json:"rating,omitempty" db:"rating"`
}
