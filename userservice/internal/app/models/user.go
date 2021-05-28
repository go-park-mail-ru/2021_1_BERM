package models

import "github.com/lib/pq"

//easyjson:json
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

//easyjson:json
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
	Executor        bool           `json:"executor,omitempty" db:"executor"`
}

//easyjson:json
type UserInfoList []UserInfo

//easyjson:json
type UserInfo struct {
	ID          uint64         `json:"id" db:"id"`
	Email       string         `json:"email" db:"email"`
	Login       string         `json:"login" db:"login"`
	NameSurname string         `json:"name_surname" db:"name_surname"`
	Password    []byte         `json:"-" db:"password"`
	About       string         `json:"about,omitempty" db:"about"`
	Specializes pq.StringArray `json:"specializes" db:"specializes,omitempty"`
	Executor    bool           `json:"executor,omitempty" db:"executor"`
	Img         string         `json:"img,omitempty" db:"img"`
	Rating      float64        `json:"rating" db:"rating"`
	ReviewCount uint64         `json:"reviews_count" db:"reviews_count"`
}

//easyjson:json
type UserBasicInfo struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email,omitempty"`
	About       string `json:"about,omitempty"`
	Executor    bool   `json:"executor"`
	Login       string `json:"login,omitempty"`
	NameSurname string `json:"name_surname,omitempty"`
}

//easyjson:json
type SuggestUsersTittleList []SuggestUsersTittle

//easyjson:json
type SuggestUsersTittle struct {
	Title string `json:"title" db:"name_surname"`
}
