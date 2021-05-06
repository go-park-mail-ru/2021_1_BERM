package error

import "github.com/pkg/errors"

var (
	ErrorNoRows     = errors.New("Error no rows.")
	ErrorDuplicate  = errors.New("Error duplicate.")
	ErrorDataSource = errors.New("Error in data source.")
)

const (
	PostgreDuplicateErrorCode = "23505"
	InternalServerErrorMsg    = "Ooops. Something went wrong!!! :("
)
