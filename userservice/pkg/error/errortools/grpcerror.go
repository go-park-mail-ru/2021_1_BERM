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
func grpcErrorListCreate() []grpcErrorInfo {
	return []grpcErrorInfo{
		{Code: customError.GRPCInternalErrorCode, Handler: func(err error) (interface{}, int){
			return map[string]interface{}{
				"message" : customError.InternalServerErrorMsg,
			}, http.StatusInternalServerError
		}},
		{Code: customError.GRPCValidationErrorCode, Handler: func(err error) (interface{}, int) {
			//FIXME проверить будет ли нормально отдаваться пользователю
			return err.Error(), http.StatusBadRequest
		}},
		{Code: customError.GRPCDuplicateErrorCode, Handler: func(err error) (interface{}, int) {
			return map[string]interface{}{
				"message": "This object already exists.",
			}, http.StatusBadRequest
		}},
		{Code: customError.GRPCNoDataErrorCode, Handler: func(err error) (interface{}, int) {
			return map[string]interface{}{
				"message": "Source not found.",
			}, http.StatusNotFound
		}},
	}
}


func grpcErrorHandle(err error)(interface{}, int, bool){
	grpcErr := &status.Status{}
	if errors.As(err, &grpcErr){
		grpsErrors := grpcErrorListCreate()
		for _, errorInfo := range grpsErrors{
			if errorInfo.Code == grpcErr.Code(){
				data, code := errorInfo.Handler(err)
				return data, code, true
			}
		}
	}
	return nil, 0, false
}

