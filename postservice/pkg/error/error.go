package error

import "github.com/pkg/errors"

var (
	ErrorNoRows          = errors.New("error no rows")
	ErrorDuplicate       = errors.New("error duplicate")
	ErrorDataSource      = errors.New("error in data source")
	ErrorUserNotExecutor = errors.New("select user not executor")
	ErrorSameID          = errors.New("executor and customer ID are the same")
)

const (
	PostgreDuplicateErrorCode = "23505"
	InternalServerErrorMsg    = "Ooops. Something went wrong!!! :("
)
