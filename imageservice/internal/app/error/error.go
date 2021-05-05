package error

import "github.com/pkg/errors"

var (
	ErrorNoRows     = errors.New("Error no rows.")
	ErrorDuplicate  = errors.New("Error duplicate.")
	ErrorDataSource = errors.New("Error in data source.")
	ErrorInvalidPassword = errors.New("Invalid password")
)

const (
	PostgreDuplicateErrorCode = "23505"
	InternalServerErrorMsg = "Ooops. Something went wrong!!! :("
)
