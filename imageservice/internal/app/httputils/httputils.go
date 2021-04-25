package httputils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"image/internal/app/logger"
	"net/http"
)

const (
	ctxKeyReqID uint8 = 1
)

func Respond(w http.ResponseWriter, requestId uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			//TODO: сделать ошибку
			s.error(w, requestId, err)
			return
		}
	}
	logger.LoggingResponse(requestId, code)
}

func RespondCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)
		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))
		Respond(w, reqID, http.StatusForbidden, map[string]interface{}{
			"error": "Invalid CSRF token",
		})
	})
}

func RespondError(w http.ResponseWriter, requestId uint64, err error) {
	//TODO: прикрутить ошибки
	httpError := &Error{}
	logger.LoggingError(requestId, err)
	if errors.As(err, &httpError) {
		Respond(w, requestId, httpError.Code, httpError.Field)
		return
	}
	Respond(w, requestId, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})

}
