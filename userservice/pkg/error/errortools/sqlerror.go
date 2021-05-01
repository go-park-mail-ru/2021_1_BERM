package errortools

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"net/http"
	customError "user/pkg/error"
)

func SqlErrorChoice(err error) error {
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

type sqlErrorInfo struct {
	Err     error
	Handler func(error) (map[string]interface{}, int)
}

func sqlErrorListCreate() []sqlErrorInfo {
	return []sqlErrorInfo{
		{Err: customError.ErrorSqlNoRows, Handler: func(err error) (map[string]interface{}, int) {
			return map[string]interface{}{
				"message": "Source not found.",
			}, http.StatusNotFound
		}},
		{Err: customError.ErrorDataSource, Handler: func(err error) (map[string]interface{}, int) {
			return map[string]interface{}{
				"message": customError.InternalServerErrorMsg,
			}, http.StatusInternalServerError
		}},
		{Err: customError.ErrorSqlDuplicate, Handler: func(err error) (map[string]interface{}, int) {
			return map[string]interface{}{
				"message": "This object already exists.",
			}, http.StatusBadRequest
		}},
	}
}

func sqlErrorHandle(err error)(map[string]interface{}, int, bool){
	errorList := sqlErrorListCreate()
	for _, errorInfo := range errorList{
		if errorInfo.Err == err {
			data, code := errorInfo.Handler(err)
			return data, code, true
		}
	}
	return nil, 0, false
}

