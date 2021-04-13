package tarantoolcache

import "errors"

var (
	NotAuthorized = errors.New("Not authorized.")
)

const (
	sessionSourceError = "Session source error."
)
