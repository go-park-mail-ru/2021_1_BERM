package httputils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"imageservice/internal/app/error/errortools"
	"imageservice/internal/app/logger"
	"imageservice/internal/app/metric"
	"net/http"
)

const (
	ctxKeyReqID uint8 = 1
)

func Respond(w http.ResponseWriter, r *http.Request, requestId uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			RespondError(w, r, requestId, err)
			return
		}
	}
	metric.CrateRequestHits(code, r)
	logger.LoggingResponse(requestId, code)
}

func RespondError(w http.ResponseWriter, r *http.Request, requestId uint64, err error) {
	logger.LoggingError(requestId, err)
	responseBody, code := errortools.ErrorHandle(err)
	metric.CrateRequestError(err)
	Respond(w, r, requestId, code, responseBody)
}

func RespondCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)

		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))
		Respond(w, r, reqID, http.StatusForbidden, map[string]interface{}{
			"error": "Invalid CSRF token",
		})
	})
}
