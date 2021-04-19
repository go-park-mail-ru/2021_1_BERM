package databases

import "github.com/tarantool/go-tarantool"

func NewTarantool(dbURL string) (*tarantool.Connection, error) {
	opts := tarantool.Opts{User: "guest"}
	db, err := tarantool.Connect(dbURL, opts)
	return db, err
}
