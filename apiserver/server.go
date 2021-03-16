package apiserver

import (
	"encoding/json"
	"errors"
	"fl_ru/model"
	"fl_ru/store"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	cookiesSalt = "ajsh468Slasdl*6%%8"
)

type server struct {
	router http.Handler
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store, config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	s.configureRouter(config)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(config *Config) {
	router := mux.NewRouter()
	router.HandleFunc("/profile", s.handleSignUp()).Methods(http.MethodPost)
	router.HandleFunc("/login", s.handleSignIn()).Methods(http.MethodPost)
	router.HandleFunc("/logout", s.handleLogout()).Methods(http.MethodDelete)
	router.HandleFunc("/profile/{id:[0-9]+}", s.authenticateUser(s.handleChangeProfile())).Methods(http.MethodPatch)
	router.HandleFunc("/order", s.authenticateUser(s.handleCreateOrder())).Methods(http.MethodPost)
	router.HandleFunc("/profile/avatar", s.authenticateUser(s.handlePutAvatar(config.ContentDir))).Methods(http.MethodPost)
	router.HandleFunc("/profile/{id:[0-9]+}", s.authenticateUser(s.handleGetProfile())).Methods(http.MethodGet)
	// router.HandleFunc("/profile/img", s.authenticateUser(s.handleGetImg())).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Origin,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With", "Accept"},
		AllowCredentials: true,
	})
	s.router = c.Handler(router)
}

func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		if len(cookies) == 0 {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) // неавторизованный пользователь
		}
		s.delCookies(cookies)
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
	}
}

// пока не используется
func (s *server) handleGetImg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		userIDCookie, _ := r.Cookie("id")
		id, _ := strconv.Atoi(userIDCookie.Value)
		u.ID = uint64(id)
		file, err := os.Open(u.ImgURL)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)

			return
		}
		avatar, err := ioutil.ReadAll(file)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)

			return
		}

		s.respond(w, http.StatusOK, avatar)
	}
}

func (s *server) handlePutAvatar(contentDir string) http.HandlerFunc {
	type Request struct {
		Img string `json:"img"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		currentDir := contentDir
		u := &model.User{}
		userIDCookie, _ := r.Cookie("id")
		id, _ := strconv.Atoi(userIDCookie.Value)
		u.ID = uint64(id)
		if r.Body == nil {
			s.error(w, http.StatusBadRequest, errors.New("No body")) // no body

			return
		}
		req := &Request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, http.StatusBadRequest, errors.New("Bad json")) //bad json

			return
		}
		pathLen := len(currentDir)
		if currentDir[pathLen-1] == '/' {
			currentDir = currentDir + userIDCookie.Value + ".base64"
		} else {
			currentDir = currentDir + "/" + userIDCookie.Value + ".base64"
		}
		file, err := os.Create(currentDir)
		if err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

			return
		}
		if _, err = file.Write([]byte(req.Img)); err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

			return
		}
		if err = file.Close(); err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

			return
		}
		u.ImgURL = currentDir
		if err = s.store.User().ChangeUser(u); err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

			return
		}
		u.Sanitize()
		s.respond(w, http.StatusCreated, u)
		defer r.Body.Close()
	}
}

func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err := s.store.User().Create(u)
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
}

func (s *server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		userIDCookie, err := r.Cookie("id")
		if err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized"))

			return
		}
		id, err := strconv.Atoi(userIDCookie.Value)
		if err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

			return
		}
		u.ID = uint64(id)
		if err := s.store.User().Find(u); err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized

			return
		}
		u.Sanitize()
		if len(u.ImgURL) != 0 {
			file, err := os.Open(u.ImgURL)
			if err != nil {
				s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) //Bad img url

				return
			}
			avatar, err := ioutil.ReadAll(file)
			if err != nil {
				s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

				return
			}
			u.ImgURL = string(avatar)
		}
		s.respond(w, http.StatusOK, u)
	}
}

func (s *server) handleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json

			return
		}
		pass := u.Password
		if err := s.store.User().FindByEmail(u); err != nil {
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
}

func (s *server) handleChangeProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, _ := r.Cookie("id")
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json

			return
		}
		id, _ := strconv.Atoi(userIDCookie.Value)
		u.ID = uint64(id)
		if err := u.BeforeCreate(); err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error"))

			return
		}
		if err := s.store.User().ChangeUser(u); err != nil {
			// некоректные данные о пользователе
			s.error(w, http.StatusBadRequest, errors.New("Incorrect user data"))

			return
		}
		u.Sanitize()
		s.respond(w, http.StatusOK, u)
	}
}

func (s *server) handleCreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, _ := r.Cookie("id")
		o := &model.Order{}
		if err := json.NewDecoder(r.Body).Decode(o); err != nil {
			s.error(w, http.StatusBadRequest, errors.New("Bad json")) //Bad json

			return
		}
		if err := o.Validate(); err != nil {
			s.error(w, http.StatusBadRequest, errors.New("Invalid data")) //Invalid data

			return
		}
		id, _ := strconv.Atoi(userIDCookie.Value)
		o.CustomerID = uint64(id)
		if err := s.store.Order().Create(o); err != nil {
			s.error(w, http.StatusInternalServerError, errors.New("Internal server error")) //500

			return
		}

		s.respond(w, http.StatusCreated, o)
	}
}

func (s *server) authenticateUser(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie("id")
		if err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized

			return
		}
		sessionID, err := r.Cookie("session")
		if err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized

			return
		}
		_, err = r.Cookie("executor")
		if err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized

			return
		}

		session := &model.Session{
			SessionID: sessionID.Value,
		}
		if err = s.store.Session().Find(session); err != nil {
			s.error(w, http.StatusUnauthorized, errors.New("Unauthorized")) //Unauthorized

			return
		}
		userIDInt, _ := strconv.Atoi(userID.Value)
		if uint64(userIDInt) != session.UserID {
			s.error(w, http.StatusForbidden, errors.New("Bad id")) // Bad id

			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) error(w http.ResponseWriter, code int, err error) {
	logrus.Error(err)
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
	session := &model.Session{
		SessionID: u.Email + time.Now().String(),
		UserID:    u.ID,
	}
	session.BeforeChange()
	if err := s.store.Session().Create(session); err != nil {
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
			Name:     "id",
			Value:    strconv.FormatUint(u.ID, 10),
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
