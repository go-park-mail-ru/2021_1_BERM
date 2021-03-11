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
	router.HandleFunc("/signup", s.handleSignUp()).Methods(http.MethodPost)
	router.HandleFunc("/signin", s.handleSignIn()).Methods(http.MethodPost)
	router.HandleFunc("/logout", s.handleLogout()).Methods(http.MethodGet)
	router.HandleFunc("/profile/change", s.authenticateUser(s.handleChangeProfile())).Methods(http.MethodPost)
	router.HandleFunc("/order", s.authenticateUser(s.handleCreateOrder())).Methods(http.MethodPost)
	router.HandleFunc("/profile/avatar", s.authenticateUser(s.handlePutAvatar(config.ContentDir))).Methods(http.MethodPost)
	router.HandleFunc("/profile", s.authenticateUser(s.handleGetProfile())).Methods(http.MethodGet)
	//router.HandleFunc("/profile/img", s.authenticateUser(s.handleGetImg())).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Origin,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "session", "id", "executor", "X-Requested-With", "Accept"},
		AllowCredentials: true,
	})
	s.router = c.Handler(router)

}

func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		s.delCookies(cookies)
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
	}
}

//func (s *server) handleGetImg() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		u := &model.User{}
//		userIdCookie, _ := r.Cookie("id")
//		id, _ := strconv.Atoi(userIdCookie.Value)
//		u.Id = uint64(id)
//		file, err := os.Open(u.ImgUrl)
//		if err != nil {
//			s.error(w, r, http.StatusBadRequest, err)
//			return
//		}
//		avatar, err := ioutil.ReadAll(file)
//		if err != nil {
//			s.error(w, r, http.StatusBadRequest, err)
//			return
//		}
//		w.Header().Set("Content-Type", "image/jpeg")
//		s.respond(w, r, http.StatusOK, avatar)
//	}
//}

func (s *server) handlePutAvatar(contentDir string) http.HandlerFunc {
	type Request struct{
		Img string `json:"img"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		currentDir := contentDir
		u := &model.User{}
		userIdCookie, _ := r.Cookie("id")
		id, _ := strconv.Atoi(userIdCookie.Value)
		u.Id = uint64(id)
		if r.Body == nil {
			s.error(w, r, http.StatusBadRequest, errors.New("No body"))
			return
		}
		req := &Request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("Bad body"))
			return
		}
		pathLen := len(currentDir)
		if currentDir[pathLen-1] == '/' {
			currentDir = currentDir + userIdCookie.Value + ".base64"
		} else {
			currentDir = currentDir + "/" + userIdCookie.Value + ".base64"
		}
		file, err := os.Create(currentDir)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if _, err = file.Write([]byte(req.Img)); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if err = file.Close(); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if err = s.store.User().ChangeUser(u); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusOK, u)
		defer r.Body.Close()
	}
}

func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := u.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := u.BeforeCreate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err := s.store.User().Create(u)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u.Sanitize()
		cookies, err := s.createCookies(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		for _, cookie := range cookies {
			http.SetCookie(w, &cookie)
		}
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		userIdCookie, _ := r.Cookie("id")
		id, _ := strconv.Atoi(userIdCookie.Value)
		u.Id = uint64(id)
		if err := s.store.User().Find(u); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u.Sanitize()
		file, err := os.Open(u.ImgUrl)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if len(u.ImgUrl) != 0 {
			avatar, err := ioutil.ReadAll(file)
			if err != nil{
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			u.ImgUrl = string(avatar)
			s.respond(w, r, http.StatusOK, u)
		}
	}
}

func (s *server) handleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		pass := u.Password
		if err := s.store.User().FindByEmail(u); err != nil {
			s.error(w, r, http.StatusConflict, err)
			return
		}
		if u.ComparePassword(pass) == false {
			s.error(w, r, http.StatusConflict, errors.New("Bad password"))
			return
		}
		u.Sanitize()
		cookies, err := s.createCookies(u)
		if err != nil {
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
	return func(w http.ResponseWriter, r *http.Request) {
		userIdCookie, _ := r.Cookie("id")
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.Atoi(userIdCookie.Value)
		u.Id = uint64(id)
		if err := u.BeforeCreate(); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if err := s.store.User().ChangeUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusAccepted, u)
	}
}

func (s *server) handleCreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdCookie, _ := r.Cookie("id")
		o := &model.Order{}
		if err := json.NewDecoder(r.Body).Decode(o); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := o.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.Atoi(userIdCookie.Value)
		o.CustomerId = uint64(id)
		if err := s.store.Order().Create(o); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, o)
	}
}

func (s *server) authenticateUser(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := r.Cookie("id")
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		sessionId, err := r.Cookie("session")
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		_, err = r.Cookie("executor")
		if err != nil {
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
		if uint64(userIdInt) != session.UserId {
			s.error(w, r, http.StatusForbidden, errors.New("Bad id"))
			return
		}
		next.ServeHTTP(w, r)
	})

}



func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	logrus.Error(err)
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)

	}
}
func (s *server) delCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Secure = true
	}
}

func (s *server) createCookies(u *model.User) ([]http.Cookie, error) {

	session := &model.Session{
		SessionId: u.Email + time.Now().String(),
		UserId:    u.Id,
	}
	session.BeforeChange()
	if err := s.store.Session().Create(session); err != nil {
		return nil, err
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    session.SessionId,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Expires:  time.Now().AddDate(0, 1, 0),
	}

	cookies := []http.Cookie{
		cookie,
		{
			Name:     "id",
			Value:    strconv.FormatUint(u.Id, 10),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Expires:  time.Now().AddDate(0, 1, 0),
		},
		{
			Name:     "executor",
			Value:    strconv.FormatBool(u.Executor),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Expires:  time.Now().AddDate(0, 1, 0),
		},
	}
	return cookies, nil
}
