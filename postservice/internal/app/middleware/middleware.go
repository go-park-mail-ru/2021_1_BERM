package middleware

import (
	"context"
	"github.com/gorilla/csrf"
	"github.com/rs/cors"
	logger2 "image/internal/app/logger"
	"math/rand"
	"net/http"
	"post/internal/app/httputils"
)



const (
	ctxKeyReqID uint8 = 1
)

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := rand.Uint64()
		logger2.LoggingRequest(reqID, r.URL, r.Method)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyReqID, reqID)))
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
		csrf.ErrorHandler(httputils.RespondCSRF()))
}
