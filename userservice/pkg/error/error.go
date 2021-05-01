package error

import "github.com/pkg/errors"

var (
	ErrorSqlNoRows    = errors.New("Error sql no rows.")
	ErrorSqlDuplicate = errors.New("Error sql duplicate.")
	ErrorDataSource   = errors.New("Error in data source.")
)

const (
	PostgreDuplicateErrorCode = "23505"
	InternalServerErrorMsg = "Ooops. Something went wrong!!! :("
)

const (
	GRPCInternalErrorCode = 1
	GRPCValidationErrorCode = 2;
	GRPCDuplicateErrorCode = 3;
	GRPCNoDataErrorCode = 4;
)


