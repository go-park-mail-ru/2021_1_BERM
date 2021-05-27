package errortools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	customError "imageservice/internal/app/error"
	"net/http"
)

const grpcErrCode codes.Code = 16

func grpcErrorHandle(err error) (interface{}, int, bool) {
	if grpcErr, ok := status.FromError(err); ok {
		if grpcErr.Code() <= grpcErrCode {
			return map[string]string{
				"message": customError.InternalServerErrorMsg,
			}, http.StatusInternalServerError, true
		}
		return err.Error(), int(grpcErr.Code()), true
	}
	return nil, 0, false
}
