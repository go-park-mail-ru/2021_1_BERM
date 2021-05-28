package httputils

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"net/http"
	"user/pkg/error/errortools"
	"user/pkg/logger"
	"user/pkg/metric"
	"user/pkg/types"
)

const (
	ctxKeyReqID types.CtxKey = 1
)

func Respond(w http.ResponseWriter, r *http.Request, requestId uint64, code int, data []byte) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(data)
		if err != nil {
			RespondError(w, r, requestId, err)
			return
		}
	}
	//metric.CrateRequestHits(code, r)
	logger.LoggingResponse(requestId, code)
}

func RespondError(w http.ResponseWriter, r *http.Request, requestId uint64, err error) {
	logger.LoggingError(requestId, err)
	responseBody, code := errortools.ErrorHandle(err)
	metric.CrateRequestError(err)

	var binBuf bytes.Buffer
	binary.Write(&binBuf, binary.BigEndian, responseBody)
	Respond(w, r, requestId, code, binBuf.Bytes())
}

func RespondCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)

		response := map[string]interface{}{
			"error": "Invalid CSRF token",
		}
		var binBuf bytes.Buffer
		binary.Write(&binBuf, binary.BigEndian, response)
		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))
		Respond(w, r, reqID, http.StatusForbidden, binBuf.Bytes())
	})
}
