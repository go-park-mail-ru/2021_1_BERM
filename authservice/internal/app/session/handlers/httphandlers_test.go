package handlers

import (
	"authorizationservice/internal/app/models"
	sessionMock "authorizationservice/internal/app/session/mock"
	"authorizationservice/pkg/metric"
	"bytes"
	"context"
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
	req, err := http.NewRequest("GET", "/login", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sessionID := "sadasdsadkmsalkdsajklda"
	req.AddCookie(&http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
	})

	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	sessionUseCaseMock.EXPECT().Get(sessionID, context.Background()).Times(1).Return(&models.Session{}, nil)

	handle := New(sessionUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CheckLogin)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestLogOut(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newUSer := models.NewUser{}
	body, _ := json.Marshal(newUSer)
	req, err := http.NewRequest("GET", "/logout", bytes.NewBuffer(body))
	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sessionID := "sadasdsadkmsalkdsajklda"
	req.AddCookie(&http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
	})

	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	handle := New(sessionUseCaseMock)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.LogOut)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}
