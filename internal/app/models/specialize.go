package models

type Specialize struct {
	ID   uint64 `db:"id"`
	Name string `json:"specialize" db:"specialize_name"`
}
