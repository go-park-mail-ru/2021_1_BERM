package errortools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	customError "post/pkg/error"
)

const grcpcErrCode codes.Code = 16

func grpcErrorHandle(err error) (interface{}, int, bool) {
	if grpcErr, ok := status.FromError(err); ok {
		if grpcErr.Code() <= grcpcErrCode {
			return map[string]string{
				"message": customError.InternalServerErrorMsg,
			}, http.StatusInternalServerError, true
		}
		return err.Error(), int(grpcErr.Code()), true
	}
	return nil, 0, false
}
