package server

import (
	"context"
	"encoding/json"
	"errors"
	"ff/configs"
	"ff/internal/app/databases"
	"ff/internal/app/models"
	server "ff/server"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
	ctxKeyReqID   ctxKey = 1
)

type server struct {
	router  http.Handler
	logger  *logrus.Logger
	useCase usecase.UseCase
}

func newServer(config *configs.Config, postgres databases.Postgres, io tarantool.Connection) *server {
	s := &server{
		router:  mux.NewRouter(),
		logger:  logrus.New(),
		useCase: useCase,
	}
	s.configureRouter(config)
	if config.LogFile == "" {
		s.logger.Out = os.Stdout
	} else {
		logFileStream, err := os.Open(config.LogFile)
		if err != nil {
			logrus.Fatal(err)
		}
		s.logger.Out = logFileStream
	}
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(config *configs.Config) {
	router := mux.NewRouter()
	router.Use(s.loggingRequest)

	csrfMiddleware := csrf.Protect(
		[]byte("very-secret-string"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Secure(config.HTTPS),
		csrf.MaxAge(900),
		csrf.Path("/"),
		csrf.ErrorHandler(s.logginCSRF()))

	router.HandleFunc("/profile", s.handleProfile).Methods(http.MethodPost)
	router.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)

	logout := router.PathPrefix("/logout").Subrouter()
	logout.Use(s.authenticateUser)
	logout.Use(csrfMiddleware)
	logout.HandleFunc("", s.handleLogout).Methods(http.MethodDelete)

	profile := router.PathPrefix("/profile").Subrouter()
	profile.Use(s.authenticateUser)
	profile.Use(csrfMiddleware)
	profile.HandleFunc("/{id:[0-9]+}", s.handleChangeProfile).Methods(http.MethodPut)
	profile.HandleFunc("/{id:[0-9]+}", s.handleGetProfile).Methods(http.MethodGet)
	profile.HandleFunc("/authorized", s.handleCheckAuthorized).Methods(http.MethodGet)
	profile.HandleFunc("/{id:[0-9]+}/specialize", s.handleAddSpecialize).Methods(http.MethodPost)
	profile.HandleFunc("/{id:[0-9]+}/specialize", s.handleDelSpecialize).Methods(http.MethodDelete)
	profile.HandleFunc("/avatar", s.handlePutAvatar).Methods(http.MethodPut)
	order := router.PathPrefix("/order").Subrouter()
	order.Use(s.authenticateUser)
	order.Use(csrfMiddleware)
	order.HandleFunc("", s.handleCreateOrder).Methods(http.MethodPost)
	order.HandleFunc("", s.handleGetActualOrder).Methods(http.MethodGet)

	//order.HandleFunc("/{id:[0-9]+}", s.handleChangeOrder).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}", s.handleGetOrder).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}/response", s.handleCreateOrderResponse).Methods(http.MethodPost)
	order.HandleFunc("/{id:[0-9]+}/response", s.handleGetAllOrderResponses).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}/response", s.handleChangeOrderResponse).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}/response", s.handleDeleteOrderResponse).Methods(http.MethodDelete)
	order.HandleFunc("/{id:[0-9]+}/select", s.handleSelectExecutor).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}/select", s.handleDeleteExecutor).Methods(http.MethodDelete)
	order.HandleFunc("/profile/{id:[0-9]+}", s.handleGetAllUserOrders).Methods(http.MethodGet)

	vacancy := router.PathPrefix("/vacancy").Subrouter()
	vacancy.Use(s.authenticateUser)
	vacancy.Use(csrfMiddleware)
	vacancy.HandleFunc("", s.handleCreateVacancy).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}", s.handleGetVacancy).Methods(http.MethodGet)
	vacancy.HandleFunc("/{id:[0-9]+}/response", s.handleCreateVacancyResponse).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}/response", s.handleGetAllVacancyResponses).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Origin,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With", "Accept", "X-Csrf-Token"},
		ExposedHeaders:   []string{"X-Csrf-Token"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
	s.router = c.Handler(router)
}

func (s *server) logginCSRF() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)
		s.logger.WithFields(logrus.Fields{
			"request_id": reqID,
			"errors":      "Invalid CSRF token",
		}).Error()
		s.respond(w, reqID, http.StatusForbidden, map[string]interface{}{
			"errors": "Invalid CSRF token",
		})
	})
}

func (s *server) loggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := rand.Uint64()
		s.logger.WithFields(logrus.Fields{
			"request_id": reqID,
			"url":        r.URL,
			"method":     r.Method,
		}).Info()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyReqID, reqID)))
	})
}


func (s *server) error(w http.ResponseWriter, requestId uint64, err error) {
	httpError := &server2.Error{}
	s.logger.WithFields(logrus.Fields{
		"errors":      err.Error(),
		"field":      httpError.Field,
		"request_id": requestId,
	}).Error()
	if errors.As(err, &httpError) {
		s.respond(w, requestId, httpError.Code, httpError.Field)
		return
	}
	s.respond(w, requestId, http.StatusInternalServerError, map[string]string{"errors": "Internal server errors"})

}


func (s *server) respond(w http.ResponseWriter, requestId uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.error(w, requestId, err)
			return
		}
	}
	s.logger.WithFields(logrus.Fields{
		"request_id": requestId,
		"reply_code": code,
	}).Info()
}

func (s *server) delCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		cookie.HttpOnly = true
	}
}

func (s *server) createCookies(u *models.User) ([]http.Cookie, error) {

	session, err := s.useCase.Session().Create(u)
	if err != nil {
		return nil, err
	}

	cookies := []http.Cookie{
		{
			Name:     "session",
			Value:    session.SessionID,
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		},
		{
			Name:     "executor",
			Value:    strconv.FormatBool(u.Executor),
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		},
	}

	return cookies, nil
}
