package httputils

import (
	"bytes"
	"encoding/binary"
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

func Respond(w http.ResponseWriter, r *http.Request, requestID uint64, code int, data []byte) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(data)
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

	var binBuf bytes.Buffer
	binary.Write(&binBuf, binary.BigEndian, responseBody)
	Respond(w, r, requestID, code, binBuf.Bytes())
}

func RespondCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)

		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))

		response := map[string]interface{}{
			"error": "Invalid CSRF token",
		}
		var binBuf bytes.Buffer
		binary.Write(&binBuf, binary.BigEndian, response)
		Respond(w, r, reqID, http.StatusForbidden, binBuf.Bytes())
	})
}
