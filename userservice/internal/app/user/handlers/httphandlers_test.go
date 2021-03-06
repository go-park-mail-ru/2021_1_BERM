package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user/internal/app/models"
	userHandlers "user/internal/app/user/handlers"
	userMock "user/internal/app/user/mock"
	"user/pkg/metric"
	"user/pkg/types"
)

const (
	ctxKeyReqID        types.CtxKey = 1
	ctxKeyStartReqTime types.CtxKey = 5
)

func TestCreateUserWithValidUrl(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := userMock.NewMockUseCase(ctrl)

	handle := userHandlers.New(mockUserUseCase)

	changeUser := &models.ChangeUser{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "Abcdasda",
		NameSurname: "Abc Def",
		Password:    "zxcvb1234$",
		About:       "ADSAdasd asd assad a dasd adsad asdas dsa dsa das da",
	}

	body, _ := json.Marshal(changeUser)

	req, err := http.NewRequest("POST", "/profile/1", bytes.NewBuffer(body))
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.ChangeProfile)
	mockUserUseCase.EXPECT().Change(*changeUser, req.Context()).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestCreateUserWithInvalidUrl(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := userMock.NewMockUseCase(ctrl)

	handle := userHandlers.New(mockUserUseCase)

	changeUser := &models.ChangeUser{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "Abcdasda",
		NameSurname: "Abc Def",
		Password:    "zxcvb1234$",
		About:       "ADSAdasd asd assad a dasd adsad asdas dsa dsa das da",
	}

	body, _ := json.Marshal(changeUser)

	req, err := http.NewRequest("POST", "/profile/asda", bytes.NewBuffer(body))

	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.ChangeProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestGetUsers(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := userMock.NewMockUseCase(ctrl)

	handle := userHandlers.New(mockUserUseCase)

	req, err := http.NewRequest("GET", "/profile/users", bytes.NewBuffer([]byte{}))

	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req.URL.Query().Add("suggest_word", "")
	mockUserUseCase.EXPECT().SuggestUsersTitle("", context.Background()).Times(1).Return([]models.SuggestUsersTittle{}, nil)
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SuggestUsers)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestSuggestErr(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := userMock.NewMockUseCase(ctrl)

	handle := userHandlers.New(mockUserUseCase)

	req, err := http.NewRequest("GET", "/profile/users", bytes.NewBuffer([]byte{}))

	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req.URL.Query().Add("suggest_word", "")
	mockUserUseCase.EXPECT().SuggestUsersTitle("", context.Background()).Times(1).Return(nil, errors.New("hunia"))
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SuggestUsers)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestGetUserInfo(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := userMock.NewMockUseCase(ctrl)
	handle := userHandlers.New(mockUserUseCase)

	req, err := http.NewRequest("GET", "/profile/1", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)

	mockUserUseCase.EXPECT().GetById(uint64(1), req.Context()).Times(1).Return(&models.UserInfo{}, nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetUserInfo)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}
