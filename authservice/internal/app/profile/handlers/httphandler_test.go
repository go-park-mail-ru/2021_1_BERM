package handlers

import (
	"authorizationservice/internal/app/models"
	profileMock "authorizationservice/internal/app/profile/mock"
	sessionMock "authorizationservice/internal/app/session/mock"
	"authorizationservice/pkg/metric"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const ctxKeyStartReqTime uint8 = 5
const ctxKeyReqID uint8 = 1

func TestRegistrationProfile(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newUSer := models.NewUser{}
	body, _ := json.Marshal(newUSer)
	req, err := http.NewRequest("POST", "/profile", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)
	profileUseCaseMock.EXPECT().Create(newUSer, context.Background()).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	sessionUseCaseMock.EXPECT().Create(uint64(1), true, context.Background()).Times(1).Return(&models.Session{}, nil)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.RegistrationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestRegistrationProfileErr(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	body, _ := json.Marshal("kek")
	req, err := http.NewRequest("POST", "/profile", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)

	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.RegistrationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestRegistrationProfileErr1(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newUSer := models.NewUser{}
	body, _ := json.Marshal(newUSer)
	req, err := http.NewRequest("POST", "/profile", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)
	profileUseCaseMock.EXPECT().Create(newUSer, context.Background()).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, sql.ErrNoRows)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.RegistrationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestRegistrationProfileErr3(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newUSer := models.NewUser{}
	body, _ := json.Marshal(newUSer)
	req, err := http.NewRequest("POST", "/profile", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)
	profileUseCaseMock.EXPECT().Create(newUSer, context.Background()).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	sessionUseCaseMock.EXPECT().Create(uint64(1), true, context.Background()).Times(1).Return(&models.Session{}, sql.ErrNoRows)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.RegistrationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestAuthorizationProfile(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newLogin := models.LoginUser{
		Password: "1qwqwdas",
		Email:    "1312322",
	}
	body, _ := json.Marshal(newLogin)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)
	profileUseCaseMock.EXPECT().Authentication(newLogin.Email, newLogin.Password, ctx).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	sessionUseCaseMock.EXPECT().Create(uint64(1), true, ctx).Times(1).Return(&models.Session{}, nil)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.AuthorisationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	metric.Destroy()
}

func TestAuthorizationProfileErr(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	body, _ := json.Marshal("mem")
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)

	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.AuthorisationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestAuthorizationProfileErr2(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newLogin := models.LoginUser{
		Password: "1qwqwdas",
		Email:    "1312322",
	}
	body, _ := json.Marshal(newLogin)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)
	profileUseCaseMock.EXPECT().Authentication(newLogin.Email, newLogin.Password, ctx).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, sql.ErrNoRows)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.AuthorisationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestAuthorizationProfile3(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newLogin := models.LoginUser{
		Password: "1qwqwdas",
		Email:    "1312322",
	}
	body, _ := json.Marshal(newLogin)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profileUseCaseMock := profileMock.NewMockUseCase(ctrl)
	profileUseCaseMock.EXPECT().Authentication(newLogin.Email, newLogin.Password, ctx).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	sessionUseCaseMock.EXPECT().Create(uint64(1), true, ctx).Times(1).Return(&models.Session{}, sql.ErrNoRows)

	handle := New(sessionUseCaseMock, profileUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.AuthorisationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}
