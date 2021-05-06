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
	reviewMock "user/internal/app/review/mock"
	"user/pkg/metric"
)

const ctxKeyStartReqTime uint8 = 5

func TestCreateReview(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	review := models.Review{
		ID:              1,
		UserId:          1,
		ToUserId:        2,
		OrderId:         1,
		Description:     "Збс делает",
		Score:           4,
		OrderName:       "Сделай что то",
		UserLogin:       "Lala@mail.ru",
		UserNameSurname: "Name surname",
	}

	body, _ := json.Marshal(review)

	req, err := http.NewRequest("POST", "/profile/review", bytes.NewBuffer(body))

	ctx := req.Context()
	reqID := uint64(2281488)
	ctx = context.WithValue(ctx, ctxKeyReqID, reqID)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mockReviewUseCase := reviewMock.NewMockUseCase(ctrl)
	mockReviewUseCase.EXPECT().Create(review, ctx).Times(1).Return(&review, nil)
	handle := New(mockReviewUseCase)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.Create)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestGetAllByUserID(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/profile/1/review", nil)
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

	mockReviewUseCase := reviewMock.NewMockUseCase(ctrl)
	mockReviewUseCase.EXPECT().GetAllReviewByUserId(uint64(1), ctx).Times(1).Return(&models.UserReviews{}, nil)

	handle := New(mockReviewUseCase)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetAllByUserId)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}
