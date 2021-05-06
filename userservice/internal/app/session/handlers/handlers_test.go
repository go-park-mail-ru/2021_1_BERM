package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user/internal/app/models"
	sessionMock "user/internal/app/session/mock"
	"user/pkg/middleware"

	"user/pkg/metric"
)

const ctxKeyStartReqTime uint8 = 5

func TestCreateSession(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	specialize := &models.Specialize{
		ID:   1,
		Name: "Govno",
	}
	body, _ := json.Marshal(specialize)

	req, err := http.NewRequest("POST", "/profile/1/specialize", bytes.NewBuffer(body))
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
	sessionID := "sadasdsad sdKFLDASD"
	req.AddCookie(&http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
	})
	mockSessionUseCase := sessionMock.NewMockUseCase(ctrl)
	mockSessionUseCase.EXPECT().Check(sessionID, ctx).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)
	handle := New(mockSessionUseCase)

	recorder := httptest.NewRecorder()
	handler := handle.CheckSession(http.NotFoundHandler())

	handler.ServeHTTP(recorder, req)
	metric.Destroy()
}

func TestLogingMiddleWare(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("POST", "/profile/1/specialize", nil)
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
	handler := middleware.LoggingRequest(http.NotFoundHandler())
	handler.ServeHTTP(recorder, req)
	metric.Destroy()
}
