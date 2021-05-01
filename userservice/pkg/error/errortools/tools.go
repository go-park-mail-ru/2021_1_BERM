package errortools

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	customError "user/pkg/error"
)

func SqlErrorHandle(err error) error {
	if err == nil {
		return nil
	}

	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		if pqErr.Code == customError.PostgreDuplicateErrorCode {
			return customError.ErrorSqlDuplicate
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return customError.ErrorSqlNoRows
	}
	return customError.ErrorDataSource
}


func GrpcErrorHandle(err error ) error{
}
