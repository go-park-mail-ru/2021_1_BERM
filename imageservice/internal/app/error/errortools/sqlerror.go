package errortools

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	customError "imageservice/internal/app/error"
	"net/http"
)

func SqlErrorChoice(err error) error {
	if err == nil {
		return nil
	}
	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		if pqErr.Code == customError.PostgreDuplicateErrorCode {
			return customError.ErrorDuplicate
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return customError.ErrorNoRows
	}
	return customError.ErrorDataSource
}

type sqlErrorInfo struct {
	Err     error
	Handler func(error) (map[string]interface{}, int)
}

func sqlErrorListCreate() []sqlErrorInfo {
	return []sqlErrorInfo{
		{Err: customError.ErrorNoRows, Handler: func(err error) (map[string]interface{}, int) {
			return map[string]interface{}{
				"message": "Source not found.",
			}, http.StatusNotFound
		}},
		{Err: customError.ErrorDataSource, Handler: func(err error) (map[string]interface{}, int) {
			return map[string]interface{}{
				"message": customError.InternalServerErrorMsg,
			}, http.StatusInternalServerError
		}},
		{Err: customError.ErrorDuplicate, Handler: func(err error) (map[string]interface{}, int) {
			return map[string]interface{}{
				"message": "This object already exists.",
			}, http.StatusBadRequest
		}},
	}
}

func sqlErrorHandle(err error) (map[string]interface{}, int, bool) {
	errorList := sqlErrorListCreate()
	for _, errorInfo := range errorList {
		if errors.Is(err, errorInfo.Err) {
			data, code := errorInfo.Handler(err)
			return data, code, true
		}
	}
	return nil, 0, false
}
