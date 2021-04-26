package httputils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"user/Error"
	"user/pkg/logger"
)

func Respond(w http.ResponseWriter, requestId uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			RespondError(w, requestId, err, 500)
			return
		}
	}
	logger.LoggingResponse(requestId, code)
}

func RespondError(w http.ResponseWriter, requestId uint64, err error, errorCode int) {
	logger.LoggingError(requestId, err)
	httpError := Error.Error{}
	if errors.As(err, &httpError) {
		Respond(w, requestId, errorCode, httpError.ErrorDescription)
		return
	}
	Respond(w, requestId, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}

func RespondCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
		if err != nil {
			RespondError(w, reqID, err, http.StatusInternalServerError)
			return
		}
		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))
		Respond(w, reqID, http.StatusForbidden, map[string]interface{}{
			"error": "Invalid CSRF token",
		})
	})
}