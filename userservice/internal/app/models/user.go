package models

import "github.com/lib/pq"

//type User struct {
//	ID              uint64         `json:"id,omitempty" db:"id"`
//	Email           string         `json:"email,omitempty" db:"email"`
//	Password        string         `json:"password,omitempty" db:"-"`
//	EncryptPassword []byte         `json:"-" db:"password"`
//	Login           string         `json:"login,omitempty" db:"login"`
//	NameSurname     string         `json:"name_surname,omitempty" db:"name_surname"`
//	Executor        bool           `json:"executor,omitempty" db:"executor"`
//	About           string         `json:"about,omitempty" db:"about"`
//	Specializes     pq.StringArray `json:"specializes,omitempty" db:"specializes,omitempty"`
//	Img             string         `json:"img,omitempty" db:"img"`
//	Rating          int            `json:"rating,omitempty" db:"rating"`
//}

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
	Rating      int32            `json:"rating,omitempty" db:"rating"`
}
