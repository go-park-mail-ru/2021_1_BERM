package middleware

import (
	"context"
	"github.com/gorilla/csrf"
	"github.com/rs/cors"
	"math/rand"
	"net/http"
	httputils2 "post/pkg/httputils"
	logger "post/pkg/logger"
	"time"
)


const (
	ctxKeyReqID uint8 = 1
	ctxKeyStartReqTime uint8 = 5
)

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := rand.Uint64()
		logger.LoggingRequest(reqID, r.URL, r.Method)
		ctx := context.WithValue(r.Context(), ctxKeyStartReqTime, time.Now())
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, ctxKeyReqID, reqID)))
	})
}
func CorsMiddleware(origin []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   origin,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With", "Accept", "X-Csrf-Token"},
		ExposedHeaders:   []string{"X-Csrf-Token"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}

func CSRFMiddleware(https bool) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte("very-secret-string"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Secure(https),
		csrf.MaxAge(900),
		csrf.Path("/"),
		csrf.ErrorHandler(httputils2.RespondCSRF()))
}
