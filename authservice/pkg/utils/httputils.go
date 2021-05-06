package utils

import (
	models2 "authorizationservice/internal/app/models"
	"authorizationservice/pkg/error/errortools"
	"authorizationservice/pkg/logger"
	"authorizationservice/pkg/metric"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	ctxKeyReqID uint8 = 1
)
func Respond(w http.ResponseWriter, r* http.Request, requestId uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			RespondError(w, r, requestId, err)
			return
		}
	}
	metric.CrateRequestHits(code, r)
	metric.CrateRequestTiming(r.Context(), r)
	logger.LoggingResponse(requestId, code)
}

func RespondError(w http.ResponseWriter, r* http.Request, requestId uint64, err error) {
	logger.LoggingError(requestId, err)
	responseBody, code := errortools.ErrorHandle(err)
	metric.CrateRequestError(err)
	Respond(w, r, requestId, code, responseBody)
}

func RespondCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)

		logger.LoggingError(reqID, errors.New("Invalid CSRF token"))
		Respond(w, r,reqID, http.StatusForbidden, map[string]interface{}{
			"error": "Invalid CSRF token",
		})
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
	for i, _ := range cookies {
		cookies[i].Expires = time.Now().AddDate(0, 0, -1)
		cookies[i].HttpOnly = true
		http.SetCookie(w, cookies[i])
	}
}
