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
	Specializes     pq.StringArray `json:"specializes" db:"specializes,omitempty"`
	Executor        bool           `db:"executor"`
}

type ChangeUser struct {
	ID              uint64         `json:"id,omitempty" db:"id"`
	Email           string         `json:"email,omitempty" db:"email"`
	Login           string         `json:"login,omitempty" db:"login"`
	NameSurname     string         `json:"name_surname,omitempty" db:"name_surname"`
	Password        string         `json:"password" db:"-"`
	NewPassword     string         `json:"new_password,omitempty" db:"-"`
	EncryptPassword []byte         `json:"-" db:"password"`
	About           string         `json:"about,omitempty" db:"about"`
	Specializes     pq.StringArray `json:"specializes,omitempty" db:"specializes,omitempty"`
	Executor        bool           `db:"executor"`
}

type UserInfo struct {
	ID          uint64         `json:"id" db:"id"`
	Email       string         `json:"email" db:"email"`
	Login       string         `json:"login" db:"login"`
	NameSurname string         `json:"name_surname" db:"name_surname"`
	Password    []byte         `json:"-" db:"password"`
	About       string         `json:"about,omitempty" db:"about"`
	Specializes pq.StringArray `json:"specializes" db:"specializes,omitempty"`
	Executor    bool           `db:"executor"`
	Img         string         `json:"img,omitempty" db:"img"`
	Rating      int32          `json:"rating,omitempty" db:"rating"`
}

type UserBasicInfo struct {
	ID       uint64
	Executor bool
}
