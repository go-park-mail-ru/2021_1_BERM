package order

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"post/internal/app/models"
	"post/internal/app/order/mock"
	"post/pkg/metric"
	"testing"
	"time"
)

const ctxKeyStartReqTime uint8 = 5

func TestCreateOrder(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	order := models.Order{
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
	}

	retOrder := &models.Order{
		ID:          1,
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
		UserLogin:   "astlok",
	}

	body, _ := json.Marshal(order)

	req, err := http.NewRequest("POST", "/api/order", bytes.NewBuffer(body))

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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CreateOrder)
	mockUseCase.EXPECT().
		Create(order, context.Background()).
		Times(1).
		Return(retOrder, nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected, _ := json.Marshal(retOrder)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())

	mockUseCase.EXPECT().
		Create(order, context.Background()).
		Times(1).
		Return(retOrder, sql.ErrNoRows)

	req, err = http.NewRequest("POST", "/api/order", bytes.NewBuffer(body))
	ctx = req.Context()
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())

	req = req.WithContext(ctx)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if status := rr2.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	var byte22 string
	byte22 = "kek"
	req, err = http.NewRequest("POST", "/api/order", bytes.NewBuffer([]byte(byte22)))
	ctx = req.Context()
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	rr2 = httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if status := rr2.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestGetActualOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	retOrder := []models.Order{
		{
			ID:          1,
			OrderName:   "Сверстать сайт",
			Category:    "Back",
			CustomerID:  1,
			Deadline:    1617004533,
			Budget:      1488,
			Description: "Pomogite sdelat API",
			UserLogin:   "astlok",
		},
	}

	req, err := http.NewRequest("GET", "/api/order", nil)

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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetActualOrder)
	mockUseCase.EXPECT().
		GetActualOrders(context.Background()).
		Times(1).
		Return(retOrder, nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(retOrder)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())
}

func TestGetActualOrderErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	retOrder := []models.Order{
		{
			ID:          1,
			OrderName:   "Сверстать сайт",
			Category:    "Back",
			CustomerID:  1,
			Deadline:    1617004533,
			Budget:      1488,
			Description: "Pomogite sdelat API",
			UserLogin:   "astlok",
		},
	}

	req, err := http.NewRequest("GET", "/api/order", nil)

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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetActualOrder)
	mockUseCase.EXPECT().
		GetActualOrders(context.Background()).
		Times(1).
		Return(retOrder, sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

}

//func Router() *mux.Router {
//	r := mux.Router()
//	r.HandleFunc("/api/order/{1}", GetRequest)
//	return r
//}
//func TestGetOrder(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUseCase := mock.NewMockUseCase(ctrl)
//
//	handle := NewHandler(mockUseCase)
//
//	retOrder := &models.Order{
//		ID:          1,
//		OrderName:   "Сверстать сайт",
//		Category:    "Back",
//		CustomerID:  1,
//		Deadline:    1617004533,
//		Budget:      1488,
//		Description: "Pomogite sdelat API",
//		UserLogin:   "astlok",
//	}
//
//	req, err := http.NewRequest("GET", "/api/order/1", nil)
//
//	ctx := req.Context()
//	var val1 uint64
//	var val2 uint64
//	val1 = 1
//	val2 = 2281488
//	ctx = context.WithValue(ctx, ctxUserID, val1)
//	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
//	req = req.WithContext(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(handle.GetOrder)
//	mockUseCase.EXPECT().
//		FindByID(1, context.Background()).
//		Times(1).
//		Return(retOrder, nil)
//
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	expected, _ := json.Marshal(retOrder)
//	expectedStr := string(expected) + "\n"
//	require.Equal(t, expectedStr, rr.Body.String())
//}
