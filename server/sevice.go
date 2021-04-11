package server

import (
	"FL_2/model"
	"FL_2/usecase"
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
	router  http.Handler
	logger  *logrus.Logger
	useCase usecase.UseCase
}

func newServer(useCase usecase.UseCase, config *Config) *server {
	s := &server{
		router:  mux.NewRouter(),
		logger:  logrus.New(),
		useCase: useCase,
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
	logout.HandleFunc("", s.handleLogout).Methods(http.MethodDelete)

	profile := router.PathPrefix("/profile").Subrouter()
	profile.Use(s.authenticateUser)
	profile.HandleFunc("/{id:[0-9]+}", s.handleChangeProfile).Methods(http.MethodPut)
	profile.HandleFunc("/{id:[0-9]+}", s.handleGetProfile).Methods(http.MethodGet)
	profile.HandleFunc("/authorized", s.handleCheckAuthorized).Methods(http.MethodGet)
	profile.HandleFunc("/{id:[0-9]+}/specialize", s.handleAddSpecialize).Methods(http.MethodPost)
	profile.HandleFunc("/{id:[0-9]+}/specialize", s.handleDelSpecialize).Methods(http.MethodDelete)
	profile.HandleFunc("/avatar", s.handlePutAvatar).Methods(http.MethodPut)
	order := router.PathPrefix("/order").Subrouter()
	order.Use(s.authenticateUser)
	order.HandleFunc("", s.handleCreateOrder).Methods(http.MethodPost)
	order.HandleFunc("", s.handleGetActualOrder).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}", s.handleChangeOrder).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}", s.handleChangeOrder).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}/response", s.handleCreateResponse).Methods(http.MethodPost)
	order.HandleFunc("/{id:[0-9]+}/response", s.handleGetAllResponses).Methods(http.MethodGet)

	vacancy := router.PathPrefix("/vacancy").Subrouter()
	vacancy.Use(s.authenticateUser)
	vacancy.HandleFunc("", s.handleCreateVacancy).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}", s.handleGetVacancy).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Origin,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With", "Accept"},
		AllowCredentials: true,
	})
	s.router = c.Handler(router)
}

func (s *server) handleCreateResponse(w http.ResponseWriter, r *http.Request) {
	response := &model.Response{}
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Order not found")) //Bad json
		return
	}
	response.OrderID = id
	response, err = s.useCase.Response().Create(*response)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad body"))
		return
	}
	s.respond(w, http.StatusCreated, response)
}

func (s *server) handleGetAllResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id")) //Bad json
		return
	}
	responses, err := s.useCase.Response().FindByID(id)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id")) //Bad json
		return
	}

	s.respond(w, http.StatusOK, responses)
}

func (s *server) handleProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json
		return
	}
	err := s.useCase.User().Create(u)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Email duplicate")) //Такой имейл уже существует
		return
	}
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
	u, err := s.useCase.User().UserVerification(u.Email, u.Password)
	if err != nil {
		s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized
		return
	}
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
		executor, err := r.Cookie("executor")
		if err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Not executor")) //Unauthorized
			return
		}
		session, err := s.useCase.Session().FindBySessionID(sessionID.Value)
		if err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized
			return
		}
		session.Executor, err = strconv.ParseBool(executor.Value)
		if err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) //Unauthorized
			return
		}
		//TODO: перенести в usecase
		session.SessionId = ""
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, session)))
	})
}

func (s *server) handleChangeProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
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
	u, err = s.useCase.User().ChangeUser(*u)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Incorrect user data"))
		return
	}
	s.respond(w, http.StatusOK, u)
}

func (s *server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id"))
		return
	}
	u, err := s.useCase.User().FindByID(id)
	if err != nil {
		if u == nil {
			s.error(w, http.StatusNotFound, errors.New("user not found"))
		} else {
			s.error(w, http.StatusInternalServerError, errors.New("InternalServerError"))
		}
		return
	}
	s.respond(w, http.StatusOK, u)
}

func (s *server) handleCheckAuthorized(w http.ResponseWriter, r *http.Request) {

	session := r.Context().Value(ctxKeySession).(*model.Session)
	s.respond(w, http.StatusOK, session)
}

func (s *server) handleAddSpecialize(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	specialize := &model.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json"))
		return
	}

	if err := s.useCase.User().AddSpecialize(specialize.Name, id); err != nil {
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
	if err := s.useCase.User().DelSpecialize(specialize.Name, userID); err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))
		return
	}
	var emptyInterface interface{}
	s.respond(w, http.StatusCreated, emptyInterface)
}

func (s *server) handlePutAvatar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := &model.User{}
	err := json.NewDecoder(r.Body).Decode(u)
	u.ID = r.Context().Value(ctxKeySession).(*model.Session).UserId
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad body"))
		return
	}

	u, err = s.useCase.Media().SetImage(u, []byte(u.Img))
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad body"))
		return
	}
	s.respond(w, http.StatusOK, u)

}

func (s *server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	o := &model.Order{}
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json
		return
	}
	o.CustomerID = id
	var err error
	o, err = s.useCase.Order().Create(*o)
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
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id"))
		return
	}
	o, err := s.useCase.Order().FindByID(id)
	if err != nil {
		s.error(w, http.StatusNotFound, errors.New("Order not found"))
		return
	}
	s.respond(w, http.StatusOK, o)
}

func (s *server) handleGetActualOrder(w http.ResponseWriter, r *http.Request) {
	o, err := s.useCase.Order().GetActualOrders()
	if err != nil {
		s.error(w, http.StatusNotFound, errors.New("Orders not found"))
		return
	}
	s.respond(w, http.StatusOK, o)
}

func (s *server) handleCreateVacancy(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	v := &model.Vacancy{
		UserId: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json
		return
	}
	var err error
	if v, err = s.useCase.Vacancy().Create(*v); err != nil {
		s.error(w, http.StatusInternalServerError, errors.New("ops"))
	}
	s.respond(w, http.StatusCreated, v)
}

func (s *server) handleGetVacancy(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, errors.New("Bad id"))
		return
	}
	v, err := s.useCase.Vacancy().FindByID(id)
	if err != nil {
		s.error(w, http.StatusNotFound, errors.New("Vacancy not found"))
		return
	}
	s.respond(w, http.StatusOK, v)
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

	session, err := s.useCase.Session().Create(u)
	if err != nil {
		return nil, err
	}

	cookies := []http.Cookie{
		{
			Name:     "session",
			Value:    session.SessionId,
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