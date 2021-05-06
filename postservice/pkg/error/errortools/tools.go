package errortools

import (
	"github.com/pkg/errors"
	"net/http"
	customError "post/pkg/error"
)

func ErrorHandle(err error) (interface{}, int) {
	if respBody, code, ok := sqlErrorHandle(err); ok {
		return respBody, code
	}
	if respBody, code, ok := validationErrorHandle(err); ok {
		return respBody, code
	}
	if respBody, code, ok := grpcErrorHandle(err); ok {
		return respBody, code
	}
	if errors.Is(err, customError.ErrorSameID) {
		return map[string]interface{}{
			"message": err.Error(),
		}, http.StatusBadRequest
	}
	if errors.Is(err, customError.ErrorUserNotExecutor) {
		return map[string]interface{}{
			"message": err.Error(),
		}, http.StatusBadRequest
	}
	return map[string]interface{}{
		"message": customError.InternalServerErrorMsg,
	}, http.StatusInternalServerError
}
