package handlers

import (
	"authorizationservice/internal/app/models"
	profileMock "authorizationservice/internal/app/profile/mock"
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
	newUSer := models.NewUser{
		ID : 1,
	}
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
	profileUseCaseMock.EXPECT().Create(newUSer, ctx).Times(1).Return(&models.UserBasicInfo{
		ID : 1,
		Executor: true,
	}, nil)
	sessionUseCaseMock := sessionMock.NewMockUseCase(ctrl)
	sessionUseCaseMock.EXPECT().Create(uint64(1), true, ctx).Times(1).Return(&models.Session{}, nil)

	handle := New(sessionUseCaseMock, profileUseCaseMock)


	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.RegistrationProfile)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	//metric.D

}

