package server

import (
	"FL_2/model"
	"FL_2/store/tarantoolcache"
	"FL_2/usecase"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
	ctxKeyReqId ctxKey = 1
)

var(
	InvalidJson = &Error{
		Err : errors.New("Invalid json."),
		Code: http.StatusBadRequest,
		Type: TypeExternal,
		Field: map[string]interface{}{
			"error" : "Invalid json",
		},
	}

	InvalidCookies = &Error{
		Err : errors.New("Invalid cookie."),
		Code: http.StatusBadRequest,
		Type: TypeExternal,
	}

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
	s.logger.Out = os.Stdout
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(config *Config) {
	router := mux.NewRouter()
	router.Use(s.loggingRequest)

	//TODO: в проде убрать secure false
	csrfMiddleware := csrf.Protect(
		[]byte("very-secret-string"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Secure(false),
		csrf.MaxAge(900),
		csrf.Path("/"))

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
	order.HandleFunc("/profile/{id:[0-9]+}", s.handleGetAllUserOrders).Methods(http.MethodDelete)

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

func (s *server) loggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := rand.Uint64()
		s.logger.WithField("Request", logrus.Fields{
			"request_id" : reqId,
			"url" : r.URL,
			"method" : r.Method,
		})
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyReqId, reqId)))
	})
}

func (s *server) handleCreateOrderResponse(w http.ResponseWriter, r *http.Request) {
	response := &model.ResponseOrder{}
	reqId := r.Context().Value(ctxKeyReqId).(uint64)

	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		s.error(w, reqId, InvalidJson) //Bad json
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, InvalidJson) //Bad json
		return
	}
	response.OrderID = id
	response, err = s.useCase.ResponseOrder().Create(*response)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w, reqId, http.StatusCreated, response)
}

func (s *server) handleGetAllOrderResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, New(err)) //Bad json
		return
	}
	responses, err := s.useCase.ResponseOrder().FindByVacancyID(id)
	if err != nil {
		s.error(w, reqId, New(err)) //Bad json
		return
	}

	s.respond(w, reqId, http.StatusOK, responses)
}

func (s *server) handleChangeOrderResponse(w http.ResponseWriter, r *http.Request) {
	response := &model.ResponseOrder{}
	reqId := r.Context().Value(ctxKeyReqId).(uint64)

	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		s.error(w, reqId, InvalidJson)
		return
	}
	params := mux.Vars(r)
	var err error
	response.OrderID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	response.UserID = r.Context().Value(ctxKeySession).(*model.Session).UserId
	responses, err := s.useCase.ResponseOrder().Change(*response)

	if err != nil {
		s.error(w, reqId, New(err))
		return
	}

	s.respond(w, reqId, http.StatusOK, responses)
}

func (s *server) handleDeleteOrderResponse(w http.ResponseWriter, r *http.Request) {
	response := &model.ResponseOrder{}
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	var err error
	response.OrderID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, New(err)) //Bad json
		return
	}
	response.UserID = r.Context().Value(ctxKeySession).(*model.Session).UserId
	err = s.useCase.ResponseOrder().Delete(*response)

	if err != nil {
		s.error(w, reqId, New(err)) //Bad json
		return
	}
	var emptyInterface interface{}

	s.respond(w, reqId, http.StatusOK, emptyInterface)
}

func (s *server) handleGetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, InvalidJson)
		return
	}
	executor, err := r.Cookie("executor")
	if err != nil {
		s.error(w, http.StatusInternalServerError, New(err))
		return
	}
	var o []model.Order
	isExecutor, err := strconv.ParseBool(executor.Value)
	if err != nil {
		s.error(w, http.StatusInternalServerError, New(err))
		return
	}
	if isExecutor {
		o, err = s.useCase.Order().FindByExecutorID(userID)
	} else {
		o, err = s.useCase.Order().FindByCustomerID(userID)
	}
	if err != nil {
		s.error(w, http.StatusNotFound, New(err))
		return
	}
	s.respond(w, reqId, http.StatusOK, o)
}

func (s *server) handleProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := &model.User{}
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, reqId, InvalidJson) //Bad json
		return
	}
	err := s.useCase.User().Create(u)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	cookies, err := s.createCookies(u)
	if err != nil {
		s.error(w, reqId, New(err)) //ошибка создания сессии
		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
	s.respond(w, reqId, http.StatusCreated, u)
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	u := &model.User{}
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, reqId,InvalidJson) //Bad json
		return
	}
	u, err := s.useCase.User().UserVerification(u.Email, u.Password)
	if err != nil {
		s.error(w, reqId, New(err)) //Unauthorized
		return
	}
	cookies, err := s.createCookies(u)
	if err != nil {
		s.error(w, reqId, New(err)) // ошибка создания сессии
		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
	s.respond(w, reqId, http.StatusOK, u)
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
		reqId := r.Context().Value(ctxKeyReqId).(uint64)
		sessionID, err := r.Cookie("session")
		cookieErr := *InvalidCookies
		if err != nil {
			cookieErr.Field = map[string]interface{}{
				"session_id" : "absent",
				"status" : "Not uthorized",
			}
			s.error(w, reqId, &cookieErr)
			return
		}
		executor, err := r.Cookie("executor")
		if err != nil {
			cookieErr.Field = map[string]interface{}{
				"executor" : "absent",
				"status" : "Not executor",
			}
			s.error(w, reqId, &cookieErr)
			return
		}
		session, err := s.useCase.Session().FindBySessionID(sessionID.Value)
		if err != nil {
			s.error(w, reqId, New(err)) //Unauthorized
			return
		}
		session.Executor, err = strconv.ParseBool(executor.Value)
		if err != nil {
			s.error(w, reqId, New(err)) //Unauthorized
			return
		}
		//TODO: перенести в usecase
		session.SessionId = ""
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, session)))
	})
}

func (s *server) handleChangeProfile(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, New(err)) //Bad json
		return
	}
	u := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, reqId, InvalidJson) //Bad json
		return
	}
	userCookieID := r.Context().Value(ctxKeySession).(*model.Session).UserId
	if userCookieID != id {
		s.error(w, reqId, New(tarantoolcache.NotAuthorized))
		return
	}
	u.ID = id
	u, err = s.useCase.User().ChangeUser(*u)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w, reqId, http.StatusOK, u)
}

func (s *server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	u, err := s.useCase.User().FindByID(id)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w, reqId, http.StatusOK, u)
}

func (s *server) handleCheckAuthorized(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	session := r.Context().Value(ctxKeySession).(*model.Session)
	s.respond(w, reqId, http.StatusOK, session)
}

func (s *server) handleAddSpecialize(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	specialize := &model.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		s.error(w, reqId, InvalidJson)
		return
	}

	if err := s.useCase.User().AddSpecialize(specialize.Name, id); err != nil {
		s.error(w, reqId, New(err))
		return
	}
	var emptyInterface interface{}
	s.respond(w, reqId, http.StatusCreated, emptyInterface)
}

func (s *server) handleDelSpecialize(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	userID := r.Context().Value(ctxKeySession).(*model.Session).UserId
	specialize := &model.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		s.error(w, reqId, InvalidJson)
		return
	}
	if err := s.useCase.User().DelSpecialize(specialize.Name, userID); err != nil {
		s.error(w, reqId, New(err))
		return
	}
	var emptyInterface interface{}
	s.respond(w, reqId, http.StatusCreated, emptyInterface)
}

func (s *server) handlePutAvatar(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	defer r.Body.Close()
	u := &model.User{}
	err := json.NewDecoder(r.Body).Decode(u)
	u.ID = r.Context().Value(ctxKeySession).(*model.Session).UserId
	if err != nil {
		s.error(w, reqId,InvalidJson)
		return
	}

	u, err = s.useCase.Media().SetImage(u, []byte(u.Img))
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w,reqId, http.StatusOK, u)

}


func (s *server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	o := &model.Order{}
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		s.error(w, reqId, InvalidJson) //Bad json
		return
	}
	o.CustomerID = id
	var err error
	o, err = s.useCase.Order().Create(*o)
	if err != nil {
		s.error(w, reqId, New(err)) //500
		return
	}
	s.respond(w, reqId, http.StatusCreated, o)
}

func (s *server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	o, err := s.useCase.Order().FindByID(id)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w, reqId, http.StatusOK, o)
}

func (s *server) handleGetActualOrder(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	o, err := s.useCase.Order().GetActualOrders()
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w, reqId, http.StatusOK, o)
}

func (s *server) handleCreateVacancy(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	id := r.Context().Value(ctxKeySession).(*model.Session).UserId
	v := &model.Vacancy{
		UserId: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.error(w, reqId, InvalidJson) //Bad json
		return
	}
	var err error
	if v, err = s.useCase.Vacancy().Create(*v); err != nil {
		s.error(w, reqId, New(err))
	}
	s.respond(w, reqId, http.StatusCreated, v)
}

func (s *server) handleGetVacancy(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, reqId, InvalidJson)
		return
	}
	v, err := s.useCase.Vacancy().FindByID(id)
	if err != nil {
		s.error(w, reqId, New(err))
		return
	}
	s.respond(w, reqId, http.StatusOK, v)
}

func (s *server) error(w http.ResponseWriter, requestId uint64,  err error) {
	httpError := &Error{}
	if errors.As(err, &httpError) {
		s.respond(w, requestId, httpError.Code, httpError.Field)
		return
	}
	s.respond(w, requestId, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	s.logger.WithField(httpError.Type, map[string]interface{}{
		"error":      err.Error(),
		"field":      httpError.Field,
		"request_id": requestId,
	}).Error()
}
func (s *server) handleCreateVacancyResponse(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	response := &model.ResponseVacancy{}
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		s.error(w, http.StatusBadRequest, InvalidJson) //Bad json
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, New(err)) //Bad json
		return
	}
	response.VacancyID = id
	response, err = s.useCase.ResponseVacancy().Create(*response)
	if err != nil {
		s.error(w, http.StatusBadRequest, New(err))
		return
	}
	s.respond(w, reqId, http.StatusCreated, response)
}

func (s *server) handleGetAllVacancyResponses (w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqId).(uint64)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		s.error(w, http.StatusBadRequest, InvalidJson) //Bad json
		return
	}
	responses, err := s.useCase.ResponseVacancy().FindByVacancyID(id)
	if err != nil {
		s.error(w, http.StatusBadRequest, New(err)) //Bad json
		return
	}

	s.respond(w,reqId,http.StatusOK, responses)
}


func (s *server) respond(w http.ResponseWriter, requestId uint64, code int, data interface{}) {
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil{
			s.error(w, requestId, err)
			return
		}
	}
	w.WriteHeader(code)
	s.logger.WithField("Reply to request", logrus.Fields{
		"request_id" : requestId,
		"reply_code" : code,
	}).Info()
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
