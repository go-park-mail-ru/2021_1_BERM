package image_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	handler "imageservice/internal/app/image/handlers"
	"imageservice/internal/app/image/mock"
	"imageservice/internal/app/metric"
	"imageservice/internal/app/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	ctxKeyReqID        uint8 = 1
	ctxKeyStartReqTime uint8 = 5
)

func TestHandlers_PutAvatar(t *testing.T) {
	metric.New()

	u := models.UserImg{
		ID:  1,
		Img: "kek",
	}

	body, _ := json.Marshal(u)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := handler.NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/profile/avatar", bytes.NewBuffer(body))

	ctx := req.Context()

	val2 := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())

	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.PutAvatar)
	mockUseCase.EXPECT().
		SetImage(u).
		Times(1).
		Return(u, nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestHandlers_PutAvatarErr(t *testing.T) {
	metric.New()

	u := models.UserImg{
		ID:  1,
		Img: "kek",
	}

	body, _ := json.Marshal(u)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := handler.NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/profile/avatar", bytes.NewBuffer(body))

	ctx := req.Context()

	val2 := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())

	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.PutAvatar)
	mockUseCase.EXPECT().
		SetImage(u).
		Times(1).
		Return(u, errors.New("err"))

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_PutAvatarErrJson(t *testing.T) {
	metric.New()

	body, _ := json.Marshal("kek")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := handler.NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/profile/avatar", bytes.NewBuffer(body))

	ctx := req.Context()
	var val2 uint64
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())

	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.PutAvatar)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}
