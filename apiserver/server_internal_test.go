package apiserver

import (
	"bytes"
	"encoding/json"
	"fl_ru/model"
	"fl_ru/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHandle_SignUp(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}

	u := model.TestUser(t)
	bCorrectUser, err := json.Marshal(u)
	assert.NoError(t, err)
	u.Email = "asdasd"
	bIncorrectUser, err := json.Marshal(u)
	assert.NoError(t, err)
	u.Email = "asdasd@gmail.com"
	u.Password = "sada"
	bIncorrectPassword, err := json.Marshal(u)
	testCases := []struct {
		Name    string
		ReqBody []byte
	}{
		{
			Name:    "CorrectRequest",
			ReqBody: bCorrectUser,
		},
		{
			Name:    "IncorrectBodyRequest",
			ReqBody: []byte("ADSSADAS"),
		},
		{
			Name:    "IncorrectEmailRequest",
			ReqBody: bIncorrectUser,
		},
		{
			Name:    "IncorrectPasswordRequest",
			ReqBody: bIncorrectPassword,
		},
	}

	handler := server.handleSignUp()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/profile", bytes.NewReader(testCase.ReqBody))
			handler.ServeHTTP(rec, req)
			var testUser model.User
			err = json.Unmarshal(rec.Body.Bytes(), &testUser)
			assert.NoError(t, err)
			switch testCase.Name {
			case "CorrectRequest":
				assert.NotEmpty(t, testUser.ID)
				assert.Empty(t, testUser.Password)
				assert.Equal(t, http.StatusCreated, rec.Code)
			case "IncorrectBodyRequest":
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			case "IncorrectEmailRequest":
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			case "IncorrectPasswordRequest":
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestHandle_SignIn(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	TestCreateCorrectUserInStore(t, server)
	u := model.TestUser(t)
	bCorrectLogin, err := json.Marshal(&model.User{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return
	}
	bIncorrectLogin, err := json.Marshal(&model.User{
		Email:    "12313",
		Password: u.Password,
	})
	if err != nil {
		return
	}
	testCases := []struct {
		Name    string
		ReqBody []byte
	}{
		{
			Name:    "CorrectRequest",
			ReqBody: bCorrectLogin,
		},
		{
			Name:    "IncorrectRequest",
			ReqBody: []byte("1231123"),
		},
		{
			Name:    "IncorrectLoginRequest",
			ReqBody: bIncorrectLogin,
		},
	}
	handler := server.handleSignIn()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader(testCase.ReqBody))
			handler.ServeHTTP(rec, req)
			var testUser model.User
			switch testCase.Name {
			case "CorrectRequest":
				err = json.Unmarshal(rec.Body.Bytes(), &testUser)
				assert.NoError(t, err)
				assert.NotEmpty(t, testUser.ID)
				assert.Empty(t, testUser.Password)
				assert.Equal(t, http.StatusOK, rec.Code)
			case "IncorrectRequest":
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			case "IncorrectLoginRequest":
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			}
		})
	}
}

func TestHandle_ChangeProfile(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	TestCreateCorrectUserInStore(t, server)
	u := model.TestUser(t)
	u.About = "1234"
	bChangeDescription, err := json.Marshal(u)
	assert.NoError(t, err)
	testCases := []struct {
		Name    string
		ReqBody []byte
	}{
		{
			Name:    "CorrectChangeDescriptionRequest",
			ReqBody: bChangeDescription,
		},
		{
			Name:    "BadRequestChangeDescriptionRequest",
			ReqBody: []byte("sdadasdsa"),
		},
	}
	handler := server.handleChangeProfile()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader(testCase.ReqBody))
			req.AddCookie(&http.Cookie{
				Name:  "id",
				Value: strconv.FormatUint(u.ID, 10),
			})
			handler.ServeHTTP(rec, req)
			var testUser model.User
			err = json.Unmarshal(rec.Body.Bytes(), &testUser)
			assert.NoError(t, err)
			switch testCase.Name {
			case "ChangeDescription":
				assert.NotEmpty(t, testUser.ID)
				assert.Empty(t, testUser.Password)
				assert.Equal(t, http.StatusAccepted, rec.Code)
				assert.Equal(t, testUser.About, u.About)
			case "BadRequestChangeDescriptionRequest":
				assert.Equal(t, rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestHandle_GetProfile(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	TestCreateCorrectUserInStore(t, server)
	u := model.TestUser(t)
	bCorrectGetReq, err := json.Marshal(u)
	assert.NoError(t, err)
	testCases := []struct {
		Name    string
		ReqBody []byte
	}{
		{
			Name:    "CorrectGetReq",
			ReqBody: bCorrectGetReq,
		},
	}
	handler := server.handleGetProfile()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/profile/10", bytes.NewReader(testCase.ReqBody))
			req.AddCookie(&http.Cookie{
				Name:  "id",
				Value: strconv.FormatUint(u.ID, 10),
			})
			handler.ServeHTTP(rec, req)
			u.Sanitize()
			var testUser model.User
			err = json.Unmarshal(rec.Body.Bytes(), &testUser)
			assert.NoError(t, err)
			switch testCase.Name {
			case "ChangeDescription":
				assert.Equal(t, u, testUser)
			}
		})
	}
}

func TestHandle_CreateOrder(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	o := model.TestOrder(t)
	bCorrectOrder, err := json.Marshal(o)
	assert.NoError(t, err)
	o.OrderName = "1"
	bBadNameOrder, err := json.Marshal(o)
	testCases := []struct {
		Name    string
		ReqBody []byte
	}{
		{
			Name:    "CorrectGetReq",
			ReqBody: bCorrectOrder,
		},
		{
			Name:    "BadReqGetReq",
			ReqBody: []byte("baasdasd"),
		},
		{
			Name:    "BadNameGetReq",
			ReqBody: bBadNameOrder,
		},
	}

	handler := server.handleCreateOrder()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/order", bytes.NewReader(testCase.ReqBody))
			req.AddCookie(&http.Cookie{
				Name:  "id",
				Value: "1",
			})
			handler.ServeHTTP(rec, req)
			var testOrder model.Order
			err = json.Unmarshal(rec.Body.Bytes(), &testOrder)
			switch testCase.Name {
			case "CorrectGetReq":
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.NoError(t, err)
				assert.NotEmpty(t, testOrder.ID)
			case "BadReqGetReq":
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			case "BadNameGetReq":
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			}
		})
	}
}

func TestCreate_Cookie(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	u := model.TestUser(t)
	cookies, err := server.createCookies(u)
	assert.NoError(t, err)
	expires := time.Now().AddDate(0, 1, 0)
	for _, cookie := range cookies {
		assert.Equal(t, cookie.Expires.Month(), expires.Month())
	}
}

func TestHandle_Authenticate(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	session := model.TestSession(t)

	err := s.Session().Create(session)
	assert.NoError(t, err)

	mw := server.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cookies := []*http.Cookie{
		{
			Name:  "id",
			Value: "1",
		},
		{
			Name:  "session",
			Value: "1",
		},
		{
			Name:  "executor",
			Value: "false",
		},
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	mw.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestHandle_AuthenticateNoSession(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	session := model.TestSession(t)

	err := s.Session().Create(session)
	assert.NoError(t, err)

	mw := server.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cookies := []*http.Cookie{
		{
			Name:  "id",
			Value: "1",
		},
		{
			Name:  "session",
			Value: "12",
		},
		{
			Name:  "executor",
			Value: "false",
		},
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	mw.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, 403)
}

func TestHandle_AuthenticateBadId(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	session := model.TestSession(t)

	err := s.Session().Create(session)
	assert.NoError(t, err)

	mw := server.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cookies := []*http.Cookie{
		{
			Name:  "session",
			Value: "1",
		},
		{
			Name:  "executor",
			Value: "false",
		},
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	mw.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestHandle_AuthenticateBadSession(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	session := model.TestSession(t)

	err := s.Session().Create(session)
	assert.NoError(t, err)

	mw := server.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cookies := []*http.Cookie{
		{
			Name:  "id",
			Value: "1",
		},
		{
			Name:  "executor",
			Value: "false",
		},
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	mw.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestHandle_AuthenticateBadExecutor(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	session := model.TestSession(t)

	err := s.Session().Create(session)
	assert.NoError(t, err)

	mw := server.authenticateUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cookies := []*http.Cookie{
		{
			Name:  "id",
			Value: "1",
		},
		{
			Name:  "session",
			Value: "1",
		},
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	mw.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
}

func TestDel_Cookie(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	u := model.TestUser(t)
	cookies, err := server.createCookies(u)
	cookiesTest := make([]*http.Cookie, len(cookies))
	for i := range cookies {
		cookiesTest[i] = &cookies[i]
	}
	assert.NoError(t, err)
	expires := time.Now().AddDate(0, 0, -1)
	server.delCookies(cookiesTest)
	for _, cookie := range cookiesTest {
		assert.Equal(t, cookie.Expires.Month(), expires.Month())
	}
}

func TestHandle_LogOut(t *testing.T) {
	s := &teststore.Store{}
	server := &server{
		store: s,
	}

	u := model.TestUser(t)
	cookies, _ := server.createCookies(u)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/logout", nil)
	handler := server.handleLogout()
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	handler.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
