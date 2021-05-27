package response

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"post/internal/app/models"
	"post/internal/app/response/mock"
	"post/pkg/metric"
	"post/pkg/types"
	"testing"
	"time"
)

const ctxKeyStartReqTime types.CtxKey = 5

func TestHandlers_CreatePostResponse(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := models.Response{
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	resResponse := models.Response{
		ID:            1,
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("POST", "/api/order/1/response", bytes.NewBuffer(body))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CreatePostResponse)
	mockUseCase.EXPECT().
		Create(response, context.Background()).
		Times(1).
		Return(&resResponse, nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	metric.Destroy()
}

func TestHandlers_CreatePostResponseErrJson(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	body := "Kek"
	req, err := http.NewRequest("POST", "/api/order/1/response", bytes.NewBuffer([]byte(body)))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CreatePostResponse)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_CreatePostResponseBadMux(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := models.Response{
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("POST", "/api/order/1/response", bytes.NewBuffer(body))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CreatePostResponse)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_CreatePostResponseErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := models.Response{
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	resResponse := models.Response{
		ID:            1,
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("POST", "/api/order/1/response", bytes.NewBuffer(body))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CreatePostResponse)
	mockUseCase.EXPECT().
		Create(response, context.Background()).
		Times(1).
		Return(&resResponse, sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_GetAllPostResponses(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := []models.Response{
		{
			PostID:        1,
			UserID:        1,
			Rate:          322,
			Text:          "TESTIY ZLO",
			Time:          228,
			UserLogin:     "astlok",
			OrderResponse: true,
		},
	}

	resResponse := []models.Response{
		{
			ID:            1,
			PostID:        1,
			UserID:        1,
			Rate:          322,
			Text:          "TESTIY ZLO",
			Time:          228,
			UserLogin:     "astlok",
			OrderResponse: true,
		},
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("GET", "/api/order/1/response", bytes.NewBuffer(body))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetAllPostResponses)
	mockUseCase.EXPECT().
		FindByPostID(uint64(1), true, false, context.Background()).
		Times(1).
		Return(resResponse, nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestHandlers_GetAllPostResponsesErrVars(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := []models.Response{
		{
			PostID:        1,
			UserID:        1,
			Rate:          322,
			Text:          "TESTIY ZLO",
			Time:          228,
			UserLogin:     "astlok",
			OrderResponse: true,
		},
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("GET", "/api/order/1/response", bytes.NewBuffer(body))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetAllPostResponses)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_GetAllPostResponsesErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := []models.Response{
		{
			PostID:        1,
			UserID:        1,
			Rate:          322,
			Text:          "TESTIY ZLO",
			Time:          228,
			UserLogin:     "astlok",
			OrderResponse: true,
		},
	}

	resResponse := []models.Response{
		{
			ID:            1,
			PostID:        1,
			UserID:        1,
			Rate:          322,
			Text:          "TESTIY ZLO",
			Time:          228,
			UserLogin:     "astlok",
			OrderResponse: true,
		},
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("GET", "/api/order/1/response", bytes.NewBuffer(body))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetAllPostResponses)
	mockUseCase.EXPECT().
		FindByPostID(uint64(1), true, false, context.Background()).
		Times(1).
		Return(resResponse, sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_ChangePostResponse(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := models.Response{
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	resResponse := models.Response{
		ID:            1,
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("PATCH", "/api/order/1/response", bytes.NewBuffer([]byte(body)))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.ChangePostResponse)
	mockUseCase.EXPECT().
		Change(response, context.Background()).
		Times(1).
		Return(&resResponse, nil)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestHandlers_ChangePostResponseErrJson(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	body := "я уже не программист, я тестей на***"
	req, err := http.NewRequest("PATCH", "/api/order/1/response", bytes.NewBuffer([]byte(body)))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.ChangePostResponse)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_ChangePostResponseErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := models.Response{
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("PATCH", "/api/order/1/response", bytes.NewBuffer([]byte(body)))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"aaaaaaa": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.ChangePostResponse)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_ChangePostResponseErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	response := models.Response{
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	resResponse := models.Response{
		ID:            1,
		PostID:        1,
		UserID:        1,
		Rate:          322,
		Text:          "TESTIY ZLO",
		Time:          228,
		UserLogin:     "astlok",
		OrderResponse: true,
	}

	body, _ := json.Marshal(response)
	req, err := http.NewRequest("PATCH", "/api/order/1/response", bytes.NewBuffer([]byte(body)))

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.ChangePostResponse)
	mockUseCase.EXPECT().
		Change(response, context.Background()).
		Times(1).
		Return(&resResponse, sql.ErrNoRows)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_DelPostResponse(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	req, err := http.NewRequest("DELETE", "/api/order/1/response", nil)

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	response := models.Response{
		PostID:          1,
		UserID:          1,
		OrderResponse:   true,
		VacancyResponse: false,
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.DelPostResponse)
	mockUseCase.EXPECT().
		Delete(response, context.Background()).
		Times(1).
		Return(nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestHandlers_DelPostResponseErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	req, err := http.NewRequest("DELETE", "/api/order/1/response", nil)

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"paaamaaagiiiiteeee": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.DelPostResponse)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestHandlers_DelPostResponseErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	req, err := http.NewRequest("DELETE", "/api/order/1/response", nil)

	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	response := models.Response{
		PostID:          1,
		UserID:          1,
		OrderResponse:   true,
		VacancyResponse: false,
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.DelPostResponse)
	mockUseCase.EXPECT().
		Delete(response, context.Background()).
		Times(1).
		Return(sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}
