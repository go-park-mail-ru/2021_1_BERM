package models

type Specialize struct {
	ID   uint64 `db:"id"`
	Name string `json:"name" db:"specialize_name"`
}

