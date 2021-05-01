package error

import "github.com/pkg/errors"

var (
	ErrorSqlNoRows    = errors.New("Error sql no rows")
	ErrorSqlDuplicate = errors.New("Error sql duplicate")
	ErrorDataSource   = errors.New("Error in data source")
	ErrorValidation   = errors.New("ValidationErrors")

)

const (
	PostgreDuplicateErrorCode = "23505"
)


