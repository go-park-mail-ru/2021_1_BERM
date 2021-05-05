package errortools

import (
	"github.com/pkg/errors"
	"net/http"
	customError "image/internal/app/error"
)

func ErrorHandle(err error) (interface{}, int) {
	if errors.Is(err, customError.ErrorInvalidPassword) {
		return map[string]string{
			"message": "Invalid password.",
		}, http.StatusBadRequest
	}
	if respBody, code, ok := sqlErrorHandle(err); ok {
		return respBody, code
	}
	if respBody, code, ok := validationErrorHandle(err); ok {
		return respBody, code
	}
	if respBody, code, ok := grpcErrorHandle(err); ok {
		return respBody, code
	}
	return map[string]interface{}{
		"message": customError.InternalServerErrorMsg,
	}, http.StatusInternalServerError
}
