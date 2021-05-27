package httputils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"post/pkg/error/errortools"
	"post/pkg/logger"
	"post/pkg/metric"
	"post/pkg/types"
)

const (
	ctxKeyReqID types.CtxKey = 1
)

func Respond(w http.ResponseWriter, r *http.Request, requestID uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			RespondError(w, r, requestID, err)
			return
		}
	}
	//metric.CrateRequestHits(code, r)
	logger.LoggingResponse(requestID, code)
}

func RespondError(w http.ResponseWriter, r *http.Request, requestID uint64, err error) {
	logger.LoggingError(requestID, err)
	responseBody, code := errortools.ErrorHandle(err)
	metric.CrateRequestError(err)
	Respond(w, r, requestID, code, responseBody)
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
