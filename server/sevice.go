package server

import (
	"FL_2/cache"
	"FL_2/model"
	"FL_2/store"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
)

type server struct {
	router http.Handler
	logger *logrus.Logger
	store  store.Store
	cache  cache.Cash
}

func newServer(store store.Store, cache cache.Cash, config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		cache:  cache,
	}
	s.configureRouter(config)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(config *Config) {
	router := mux.NewRouter()
	router.HandleFunc("/profile", s.handleProfile).Methods(http.MethodPost)
	router.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)

	logout := router.PathPrefix("/logout").Subrouter()
	logout.Use(s.authenticateUser)
	logout.HandleFunc("/logout", s.handleLogout).Methods(http.MethodDelete)

	profile := router.PathPrefix("/profile").Subrouter()
	profile.Use(s.authenticateUser)
	profile.HandleFunc("/{id:[0-9]+}", s.handleChangeProfile).Methods(http.MethodPut)
	profile.HandleFunc("/{id:[0-9]+}", s.handleGetProfile).Methods(http.MethodGet)
	profile.HandleFunc("/authorized", s.handleCheckAuthorized).Methods(http.MethodGet)
	profile.HandleFunc("/{id:[0-9]+}/specialize", s.handleAddSpecialize).Methods(http.MethodPost)
	profile.HandleFunc("/{id:[0-9]+}/specialize", s.handleDelSpecialize).Methods(http.MethodDelete)

	order := router.PathPrefix("/order").Subrouter()
	order.Use(s.authenticateUser)
	order.HandleFunc("/", s.handleCreateOrder).Methods(http.MethodPost)
	order.HandleFunc("/{id:[0-9]+}", s.handleChangeOrder).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}", s.handleChangeOrder).Methods(http.MethodGet)
	c := cors.New(cors.Options{
		AllowedOrigins:   config.Origin,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With", "Accept"},
		AllowCredentials: true,
	})
	s.router = c.Handler(router)
}

func (s *server) handleProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json

		return
	}
	if err := u.Validate(); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Invalid data")) //Invalid data

		return
	}
	if err := u.BeforeCreate(); err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) //Ошибка в закодировании пароля

		return
	}

	var err error

	u.ID, err = s.store.User().Create(*u)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Email duplicate")) //Такой имейл уже существует

		return
	}
	u.Sanitize()
	cookies, err := s.createCookies(u)
	if err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) //ошибка создания сессии

		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
	s.respond(w, http.StatusCreated, u)
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	u := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json

		return
	}
	pass := u.Password
	var err error
	u, err = s.store.User().FindByEmail(u.Email)
	if err != nil {
		s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized

		return
	}
	if !u.ComparePassword(pass) {
		s.error(w, http.StatusUnauthorized, errors.New("Bad password")) //bad paswd

		return
	}
	u.Sanitize()
	cookies, err := s.createCookies(u)
	if err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) // ошибка создания сессии

		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
	s.respond(w, http.StatusOK, u)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	s.delCookies(cookies)
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionID, err := r.Cookie("session")
		if err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized
			return
		}

		session := &model.Session{
			SessionId: sessionID.Value,
		}

		if err = s.cache.Session().Find(session); err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, session)))
	})
}

func (s *server) handleChangeProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err:= strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id")) //Bad json
		return
	}
	u := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json
		return
	}
	userCookieID := r.Context().Value(ctxKeySession).(*model.Session).UserId

	if userCookieID != id {
		s.error(w, http.StatusBadRequest, errors.New("No right to modify")) //Bad json
		return
	}
	u.ID = id
	if err := u.BeforeCreate(); err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}

	u, err = s.store.User().ChangeUser(*u)
	if err != nil {
		// некоректные данные о пользователе
		s.error(w, http.StatusBadRequest, errors.New("Incorrect user data"))
		return
	}
	u.Sanitize()
	s.respond(w, http.StatusOK, u)
}

func (s *server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err:= strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id"))
		return
	}
	u := &model.User{}
	u, err = s.store.User().FindByID(id)
	if err != nil {
		s.error(w, http.StatusNotFound, errors.New("user not found"))
		return
	}
	u.Sanitize()
	s.respond(w, http.StatusOK, u)
}

func (s *server) handleCheckAuthorized(w http.ResponseWriter, r *http.Request) {
	u := &model.User{
		ID: r.Context().Value(ctxKeySession).(*model.Session).UserId,
	}

	s.respond(w, http.StatusOK, u)
}

func (s *server) handleAddSpecialize(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxKeySession).(*model.Session).UserId
	specialize := &model.Specialize{}

	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json"))
		return
	}

	if err := s.store.User().AddSpecialize(specialize.Name, userID); err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}

	var emptyInterface interface{}
	s.respond(w, http.StatusCreated, emptyInterface)
}

func (s *server) handleDelSpecialize(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxKeySession).(*model.Session).UserId
	specialize := &model.Specialize{}

	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json"))
		return
	}

	if err := s.store.User().DelSpecialize(specialize.Name, userID); err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}

	var emptyInterface interface{}
	s.respond(w, http.StatusCreated, emptyInterface)
}

func (s *server) handleAvatar(w http.ResponseWriter, r *http.Request) {

}

func (s *server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	o := &model.Order{}
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json
		return
	}
	if err := o.Validate(); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Invalid data")) //Invalid data
		return
	}
	o.CustomerID = id
	var err error
	o.ID, err = s.store.Order().Create(*o)
	if err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) //500
		return
	}

	s.respond(w, http.StatusCreated, o)
}

func (s *server) handleChangeOrder(w http.ResponseWriter, r *http.Request) {

}

func (s *server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err:= strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id"))
		return
	}
	o := &model.Order{
		ID: id,
	}
	o, err = s.store.Order().FindByID(o.ID)
	if err != nil {
		s.error(w, http.StatusNotFound, errors.New("Order not found"))
		return
	}
	s.respond(w, http.StatusOK, o)
}

func (s *server) error(w http.ResponseWriter, code int, err error) {
	s.logger.Error(err)
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *server) delCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		cookie.HttpOnly = true
	}
}

func (s *server) createCookies(u *model.User) ([]http.Cookie, error) {

	session := &model.Session{}
	session.UserId = u.ID
	session.SessionId = u.Email + time.Now().String()
	session.BeforeChange()
	if err := s.cache.Session().Create(session); err != nil {
		return nil, err
	}

	cookies := []http.Cookie{
		{
			Name:     "session",
			Value:    session.SessionId,
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
		},
	}

	return cookies, nil
}
