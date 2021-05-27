package errortools

import (
	customError "authorizationservice/pkg/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const grcpcErrCode codes.Code = 16

func GrpcErrorHandle(err error) (interface{}, int, bool) {
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
