package models

import "github.com/lib/pq"

type User struct {
	ID          uint64         `json:"id,omitempty" copier:"-"`
	Email       string         `json:"email,omitempty"`
	Login       string         `json:"login,omitempty"`
	NameSurname string         `json:"name_surname,omitempty"`
	About       string         `json:"about,omitempty"`
	Specializes pq.StringArray `json:"specializes,omitempty"`
	Executor    bool           `json:"executor,omitempty"`
	Img         string         `json:"img,omitempty"`
	Rating      uint64         `json:"rating,omitempty"`
}
