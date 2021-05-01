package errortools

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	customError "user/pkg/error"
)

type grpcErrorInfo struct {
	Code codes.Code
	Handler func(error) (interface{}, int)
}



func grpcErrorHandle(err error)(interface{}, int, bool){
	grpcErr := &status.Status{}
	if errors.As(err, &grpcErr){
		if grpcErr.Code() <= 16{
			return map[string]string{
				"message" : customError.InternalServerErrorMsg,
			}, http.StatusInternalServerError, true
		}
		return err.Error(), int(grpcErr.Code()), true
	}
	return nil, 0, false
}

