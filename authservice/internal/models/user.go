package models

import "github.com/lib/pq"

type NewUser struct {
	ID              uint64         `json:"id,omitempty" db:"id"`
	Email           string         `json:"email" db:"email"`
	Login           string         `json:"login" db:"login"`
	NameSurname     string         `json:"name_surname" db:"name_surname"`
	Password        string         `json:"password" db:"-"`
	EncryptPassword []byte         `json:"-" db:"password"`
	About           string         `json:"about,omitempty" db:"about"`
	Specializes     pq.StringArray `json:"specializes,omitempty" db:"specializes,omitempty"`
	Executor        bool           `db:"executor"`
}

type LoginUser struct {
	Email    string
	Password string
}

type UserBasicInfo struct {
	ID       uint64 `json:"id"`
	Executor bool   `json:"executor"`
}
