package handler

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
	specializeMock "user/internal/app/specialize/mock"
	userMock "user/internal/app/user/mock"
	"user/pkg/metric"
)

const ctxKeyStartReqTime uint8 = 5

func TestCreateSpecialize(t *testing.T) {
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

	mockUserUseCase := userMock.NewMockUseCase(ctrl)
	mockUserUseCase.EXPECT().GetById(uint64(1), ctx).Times(1).Return(&models.UserInfo{}, nil)

	mockSpecializeUseCase := specializeMock.NewMockUseCase(ctrl)
	mockSpecializeUseCase.EXPECT().AssociateWithUser(uint64(1), specialize.Name, ctx).Times(1).Return(nil)
	handle := New(mockSpecializeUseCase, mockUserUseCase)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.Create)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestRemoveSpecialize(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	specialize := &models.Specialize{
		ID:   1,
		Name: "Govno",
	}
	body, _ := json.Marshal(specialize)

	req, err := http.NewRequest("DELETE", "/profile/1/specialize", bytes.NewBuffer(body))
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

	mockUserUseCase := userMock.NewMockUseCase(ctrl)
	mockSpecializeUseCase := specializeMock.NewMockUseCase(ctrl)
	mockSpecializeUseCase.EXPECT().Remove(uint64(1), specialize.Name, ctx).Times(1).Return(nil)
	handle := New(mockSpecializeUseCase, mockUserUseCase)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.Remove)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}
