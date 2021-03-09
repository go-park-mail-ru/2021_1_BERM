package apiserver

import (
	"encoding/json"
	"errors"
	"fl_ru/model"
	"fl_ru/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)


const(
	cookiesSalt = "ajsh468Slasdl*6%%8"
)


type server struct{
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store: store,
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP (w http.ResponseWriter, r *http.Request){
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(){
	s.router.HandleFunc("/signup",  s.handleSignUp()).Methods("POST")
	s.router.HandleFunc("/signin",  s.handleSignIn()).Methods("POST")
	s.router.HandleFunc("/profile/change",  s.authenticateUser(s.handleChangeProfile())).Methods("POST")
	s.router.HandleFunc("/order", s.authenticateUser(s.handleCreateOrder())).Methods("POST")
	s.router.Use(s.setOptions)
}

func (s *server) setOptions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			setupDifficultResponse(&w, r)
			s.respond(w, r, http.StatusOK, nil)
			return
		}
		setupSimpleResponse(&w, r)
		next.ServeHTTP(w, r)
	})
}

func (s *server) handleSignUp() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u) ;err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := u.Validate(); err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := u.BeforeCreate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err := s.store.User().Create(u)
		if  err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u.Sanitize()
		cookies, err := s.createCookies(u)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		for _, cookie := range cookies {
			http.SetCookie(w, &cookie)
		}
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSignIn() http.HandlerFunc {
	type Request struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request){
		request := &Request{}
		if err := json.NewDecoder(r.Body).Decode(request) ;err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
			Email: request.Email,
			Password: request.Password,
		}
		if err := s.store.User().FindByEmail(u); err != nil {
			s.error(w, r, http.StatusConflict, err)
			return
		}
		if u.ComparePassword(request.Password) == false{
			s.error(w, r, http.StatusConflict, errors.New("Bad password"))
			return
		}
		u.Sanitize()
		cookies, err := s.createCookies(u)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		for _, cookie := range cookies {
			http.SetCookie(w, &cookie)
		}
		s.respond(w, r, http.StatusAccepted, u)
	}
}

func (s *server) handleChangeProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		userIdCookie, _ := r.Cookie("id")
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u) ;err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.Atoi(userIdCookie.Value)
		u.Id = uint64(id)
		if err := u.BeforeCreate(); err != nil{
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if err := s.store.User().ChangeUser(u); err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusAccepted, u)
	}
}

func (s *server) handleCreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		userIdCookie, _ := r.Cookie("id")
		o := &model.Order{}
		if err := json.NewDecoder(r.Body).Decode(o) ;err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := o.Validate() ;err != nil{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.Atoi(userIdCookie.Value)
		o.CustomerId = uint64(id)
		if err := s.store.Order().Create(o); err != nil{
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, o)
	}
}


func (s *server) authenticateUser(next http.Handler) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := r.Cookie("id")
		if err != nil{
			s.error(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		sessionId, err := r.Cookie("session")
		if err != nil{
			s.error(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		_, err = r.Cookie("executor")
		if err != nil{
			s.error(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		session := &model.Session{
			SessionId: sessionId.Value,
		}
		if err = s.store.Session().Find(session); err != nil {
			s.error(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		userIdInt, _ := strconv.Atoi(userId.Value)
		if uint64(userIdInt) !=  session.UserId {
			s.error(w, r, http.StatusUnauthorized, errors.New("Bad id"))
			return
		}
		next.ServeHTTP(w, r)
	})

}



func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error){
	s.respond(w, r, code, map[string]string{"error" : err.Error()})
}

func (s* server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}){
	w.WriteHeader(code)
	if data != nil{
		_ = json.NewEncoder(w).Encode(data)

	}
}

func (s *server)createCookies(u *model.User) ([]http.Cookie, error){

	session := &model.Session{
		SessionId: u.Email + time.Now().String(),
		UserId: u.Id,
	}
	session.BeforeChange()
	if err := s.store.Session().Create(session); err != nil{
		return nil, err
	}
	cookie := http.Cookie{
		Name: "session",
		Value: session.SessionId,
	}
	cookies := []http.Cookie{
		cookie,
		{
			Name:  "id",
			Value: strconv.FormatUint(u.Id, 10),
		},
		{
			Name: "executor",
			Value: strconv.FormatBool(u.Executor),
		},
	}
	return cookies, nil
}
