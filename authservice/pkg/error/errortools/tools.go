package errortools

import (
	customError "authorizationservice/pkg/error"
	"net/http"
)

func ErrorHandle(err error) (interface{}, int) {
	if respBody, code, ok := SqlErrorHandle(err); ok {
		return respBody, code
	}
	if respBody, code, ok := ValidationErrorHandle(err); ok {
		return respBody, code
	}
	if respBody, code, ok := GrpcErrorHandle(err); ok {
		return respBody, code
	}
	return map[string]interface{}{
		"message": customError.InternalServerErrorMsg,
	}, http.StatusInternalServerError
}
