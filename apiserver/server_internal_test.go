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

func TestHandle_SignUp(t *testing.T){

	s := &teststore.Store{}
	server := &server{
		store: s,
	}


	u := model.TestUser(t)
	bCorrectUser, err := json.Marshal(u)
	assert.NoError(t, err)
	u.Email = "asdasd"
	bIncorrectUser , err := json.Marshal(u)
	assert.NoError(t, err)
	u.Email = "asdasd@gmail.com"
	u.Password = "sada"
	bIncorrectPassword , err := json.Marshal(u)
	testCases := []struct{
		Name    string
		ReqBody []byte
	}{
		{
			Name :   "CorrectRequest",
			ReqBody: bCorrectUser,
		},
		{
			Name :   "IncorrectBodyRequest",
			ReqBody: []byte("ADSSADAS"),
		},
		{
			Name :   "IncorrectEmailRequest",
			ReqBody: bIncorrectUser,
		},
		{
			Name :   "IncorrectPasswordRequest",
			ReqBody: bIncorrectPassword,
		},
	}

	handler := http.HandlerFunc(server.handleSignUp())

	for _, testCase := range testCases{
		t.Run(testCase.Name, func(t *testing.T){
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/signup", bytes.NewReader(testCase.ReqBody))
			handler.ServeHTTP(rec, req)
			var testUser model.User
			err = json.Unmarshal(rec.Body.Bytes(), &testUser)
			assert.NoError(t, err)
			switch testCase.Name {
			case "CorrectRequest":
				assert.NotEmpty(t, testUser.Id)
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


func TestHandle_SignIn(t *testing.T){
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	TestCreateCorrectUserInStore(t, server)
	u := model.TestUser(t)
	bCorrectLogin, err := json.Marshal(&model.User{
		Email: u.Email,
		Password: u.Password,
	})
	testCases := []struct{
		Name    string
		ReqBody []byte
	}{
		{
			Name :   "CorrectRequest",
			ReqBody: bCorrectLogin,
		},
	}
	handler := http.HandlerFunc(server.handleSignIn())
	for _, testCase := range testCases{
		t.Run(testCase.Name, func(t *testing.T){
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/signin", bytes.NewReader(testCase.ReqBody))
			handler.ServeHTTP(rec, req)
			var testUser model.User
			err = json.Unmarshal(rec.Body.Bytes(), &testUser)
			assert.NoError(t, err)
			switch testCase.Name {
			case "CorrectRequest":
				assert.NotEmpty(t, testUser.Id)
				assert.Empty(t, testUser.Password)
				assert.Equal(t, http.StatusAccepted, rec.Code)
			}
		})
	}
}

func TestHandle_ChangeProfile(t *testing.T){
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	TestCreateCorrectUserInStore(t, server)
	u := model.TestUser(t)
	u.Description = "1234"
	bChangeDescription, err := json.Marshal(u)
	assert.NoError(t, err)
	testCases := []struct{
		Name    string
		ReqBody []byte
	}{
		{
			Name :   "ChangeDescription",
			ReqBody: bChangeDescription,
		},
	}
	handler := http.HandlerFunc(server.handleChangeProfile())
	for _, testCase := range testCases{
		t.Run(testCase.Name, func(t *testing.T){
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/signin", bytes.NewReader(testCase.ReqBody))
			req.AddCookie(&http.Cookie{
				Name:  "id",
				Value: strconv.FormatUint(u.Id, 10),
			})
			handler.ServeHTTP(rec, req)
			var testUser model.User
			err = json.Unmarshal(rec.Body.Bytes(), &testUser)
			assert.NoError(t, err)
			switch testCase.Name {
			case "ChangeDescription":
				assert.NotEmpty(t, testUser.Id)
				assert.Empty(t, testUser.Password)
				assert.Equal(t, http.StatusAccepted, rec.Code)
				assert.Equal(t, testUser.Description, u.Description)
			}
		})
	}
}

func TestHandle_GetProfile(t *testing.T){
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	TestCreateCorrectUserInStore(t, server)
	u := model.TestUser(t)
	bCorrectGetReq, err := json.Marshal(u)
	assert.NoError(t, err)
	testCases := []struct{
		Name    string
		ReqBody []byte
	}{
		{
			Name :   "CorrectGetReq",
			ReqBody: bCorrectGetReq,
		},
	}
	handler := http.HandlerFunc(server.handleGetProfile())
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/profile", bytes.NewReader(testCase.ReqBody))
			req.AddCookie(&http.Cookie{
				Name:  "id",
				Value: strconv.FormatUint(u.Id, 10),
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

func TestHandle_CreateOrder(t *testing.T){
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	o := model.TestOrder(t)
	bCorrectOrder, err := json.Marshal(o)
	assert.NoError(t, err)
	testCases := []struct{
		Name    string
		ReqBody []byte
	}{
		{
			Name :   "CorrectGetReq",
			ReqBody: bCorrectOrder,
		},
	}


	handler := http.HandlerFunc(server.handleCreateOrder())

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/order", bytes.NewReader(testCase.ReqBody))
			req.AddCookie(&http.Cookie{
				Name: "id",
				Value: "1",
			})
			handler.ServeHTTP(rec, req)
			var testOrder model.Order
			err = json.Unmarshal(rec.Body.Bytes(), &testOrder)

			assert.Equal(t, http.StatusAccepted, rec.Code)
			assert.NoError(t, err)
			assert.NotEmpty(t, testOrder.Id)
		})
	}
}

func TestCreate_Cookie(t *testing.T){
	s := &teststore.Store{}
	server := &server{
		store: s,
	}
	u := model.TestUser(t)
	cookies, err := server.createCookies(u)
	assert.NoError(t, err)
	expires := time.Now().AddDate(0, 1, 0)
	for _, cookie := range cookies{
		assert.Equal(t, cookie.Expires.Month(), expires.Month())
	}
}

func TestDel_Cookie(t *testing.T){
	//s := &teststore.Store{}
	//server := &server{
	//	store: s,
	//}
	//u := model.TestUser(t)
	//cookies, err := server.createCookies(u)
	//assert.NoError(t, err)

}