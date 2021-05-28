package utils

import (
	models2 "authorizationservice/internal/app/models"
	"authorizationservice/pkg/error/errortools"
	"authorizationservice/pkg/logger"
	"authorizationservice/pkg/metric"
	"authorizationservice/pkg/types"
	"bytes"
	"encoding/binary"
	"errors"
	"net/http"
	"time"
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

		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))

		response := map[string]interface{}{
			"error": "Invalid CSRF token",
		}

		var binBuf bytes.Buffer
		binary.Write(&binBuf, binary.BigEndian, response)
		Respond(w, r, reqID, http.StatusForbidden, binBuf.Bytes())
	})
}

func CreateCookie(session *models2.Session, w http.ResponseWriter) {
	cookies := []http.Cookie{
		{
			Name:     "sessionID",
			Value:    session.SessionID,
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		},
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func RemoveCookies(cookies []*http.Cookie, w http.ResponseWriter) {
	for i := range cookies {
		cookies[i].Expires = time.Now().AddDate(0, 0, -1)
		cookies[i].HttpOnly = true
		http.SetCookie(w, cookies[i])
	}
}
