package order

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
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
	metric.Destroy()
}

func TestGetActualOrder(t *testing.T) {
	metric.New()

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
	metric.Destroy()
}

//func TestGetActualOrderErr(t *testing.T) {
//	metric.New()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUseCase := mock.NewMockUseCase(ctrl)
//
//	handle := NewHandler(mockUseCase)
//
//	retOrder := []models.Order{
//		{
//			ID:          1,
//			OrderName:   "Сверстать сайт",
//			Category:    "Back",
//			CustomerID:  1,
//			Deadline:    1617004533,
//			Budget:      1488,
//			Description: "Pomogite sdelat API",
//			UserLogin:   "astlok",
//		},
//	}
//
//	req, err := http.NewRequest("GET", "/api/order?search_str=kek&from=1&to=2&desc=false&category=mem&limit=1&offset=2", nil)
//
//	ctx := req.Context()
//	var val1 uint64
//	var val2 uint64
//	val1 = 1
//	val2 = 2281488
//	param := make(map[string]interface{})
//	param["search_str"] = "kek"
//	param["from"] = "1"
//	param["to"] = "2"
//	param["desc"] = "false"
//	param["category"] = "mem"
//	param["limit"] = "1"
//	param["offset"] = "2"
//
//	ctx = context.WithValue(ctx, ctxUserID, val1)
//	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
//	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
//	ctx = context.WithValue(ctx, ctxQueryParams, param)
//	req = req.WithContext(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(handle.GetActualOrder)
//	mockUseCase.EXPECT().
//		GetActualOrders(ctx).
//		Times(1).
//		Return(retOrder, sql.ErrNoRows)
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusInternalServerError {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusInternalServerError)
//	}
//	metric.Destroy()
//}

func TestGetOrder(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

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

	req, err := http.NewRequest("GET", "/api/order/1", nil)

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
	handler := http.HandlerFunc(handle.GetOrder)
	mockUseCase.EXPECT().
		FindByID(uint64(1), context.Background()).
		Times(1).
		Return(retOrder, nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(retOrder)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())

	req1, err := http.NewRequest("GET", "/api/order/1", nil)

	ctx2 := req1.Context()
	val1 = 1
	val2 = 2281488
	ctx2 = context.WithValue(ctx2, ctxUserID, val1)
	ctx2 = context.WithValue(ctx2, ctxKeyReqID, val2)
	ctx2 = context.WithValue(ctx2, ctxKeyStartReqTime, time.Now())

	req1 = req1.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	vars = map[string]string{
		"kek": "1",
	}
	req1 = mux.SetURLVars(req1, vars)
	handler.ServeHTTP(rr, req1)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	metric.Destroy()
}

func TestGetOrderErr(t *testing.T) {

	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

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
	req, err := http.NewRequest("GET", "/api/order/1", nil)

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
	handler := http.HandlerFunc(handle.GetOrder)
	mockUseCase.EXPECT().
		FindByID(uint64(1), context.Background()).
		Times(1).
		Return(retOrder, sql.ErrNoRows)

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

func TestChangeOrder(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	order := models.Order{
		ID:          1,
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
	req, err := http.NewRequest("GET", "/api/order/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.ChangeOrder)
	mockUseCase.EXPECT().
		ChangeOrder(order, context.Background()).
		Times(1).
		Return(*retOrder, nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(retOrder)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())
	metric.Destroy()
}

func TestChangeOrderBadJson(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	byte22 := "kek"
	req, err := http.NewRequest("GET", "/api/order/1", bytes.NewBuffer([]byte(byte22)))

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
	handler := http.HandlerFunc(handle.ChangeOrder)
	//mockUseCase.EXPECT().
	//	ChangeOrder(order, context.Background()).
	//	Times(1).
	//	Return(*retOrder, nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	//expected, _ := json.Marshal(retOrder)
	//expectedStr := string(expected) + "\n"
	//require.Equal(t, expectedStr, rr.Body.String())
	metric.Destroy()
}

func TestChangeOrderErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	order := models.Order{
		ID:          1,
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
	req, err := http.NewRequest("GET", "/api/order/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.ChangeOrder)
	mockUseCase.EXPECT().
		ChangeOrder(order, context.Background()).
		Times(1).
		Return(*retOrder, sql.ErrNoRows)

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

func TestChangeOrderErrParse(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	order := models.Order{
		ID:          1,
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
	}
	body, _ := json.Marshal(order)
	req, err := http.NewRequest("GET", "/api/order/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.ChangeOrder)

	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}

func TestDeleteOrder(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/order/1", nil)

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
	handler := http.HandlerFunc(handle.DeleteOrder)
	mockUseCase.EXPECT().
		DeleteOrder(uint64(1), context.Background()).
		Times(1).
		Return(nil)

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

func TestDeleteOrderErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/order/1", nil)

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
	handler := http.HandlerFunc(handle.DeleteOrder)

	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestDeleteOrderErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/order/1", nil)

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
	handler := http.HandlerFunc(handle.DeleteOrder)
	mockUseCase.EXPECT().
		DeleteOrder(uint64(1), context.Background()).
		Times(1).
		Return(sql.ErrNoRows)

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

func TestSelectEx(t *testing.T) {
	metric.New()

	order := models.Order{
		ID:         1,
		ExecutorID: 1,
	}

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/order/1/select", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.SelectExecutor)
	mockUseCase.EXPECT().
		SelectExecutor(order, context.Background()).
		Times(1).
		Return(nil)

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

func TestSelectExErrJson(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)
	byte22 := "kek"
	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/order/1/select", bytes.NewBuffer([]byte(byte22)))

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
	handler := http.HandlerFunc(handle.SelectExecutor)
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

func TestSelectBarVar(t *testing.T) {
	metric.New()

	order := models.Order{
		ID:         1,
		ExecutorID: 1,
	}

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/order/1/select", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.SelectExecutor)

	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestSelectExErr(t *testing.T) {
	metric.New()

	order := models.Order{
		ID:         1,
		ExecutorID: 1,
	}

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/order/1/select", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.SelectExecutor)
	mockUseCase.EXPECT().
		SelectExecutor(order, context.Background()).
		Times(1).
		Return(sql.ErrNoRows)

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

func TestDeleteExecutor(t *testing.T) {
	metric.New()

	order := models.Order{
		ID: 1,
	}

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/order/1/select", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.DeleteExecutor)
	mockUseCase.EXPECT().
		DeleteExecutor(order, context.Background()).
		Times(1).
		Return(nil)

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

func TestDeleteExecutorErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/order/1/select", nil)

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
	handler := http.HandlerFunc(handle.DeleteExecutor)

	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestDeleteExecutorErr(t *testing.T) {
	metric.New()

	order := models.Order{
		ID: 1,
	}

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/order/1/select", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.DeleteExecutor)
	mockUseCase.EXPECT().
		DeleteExecutor(order, context.Background()).
		Times(1).
		Return(sql.ErrNoRows)

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

func TestGetAllOrders(t *testing.T) {
	metric.New()

	order := models.Order{
		ID: 1,
	}

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

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/order/profile/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.GetAllUserOrders)
	mockUseCase.EXPECT().
		FindByUserID(uint64(1), context.Background()).
		Times(1).
		Return(retOrder, nil)

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

func TestGetAllOrdersErrVar(t *testing.T) {
	metric.New()

	order := models.Order{
		ID: 1,
	}

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/order/profile/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.GetAllUserOrders)

	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestGetAllOrdersErr(t *testing.T) {
	metric.New()

	order := models.Order{
		ID: 1,
	}

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

	body, _ := json.Marshal(order)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/order/profile/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.GetAllUserOrders)
	mockUseCase.EXPECT().
		FindByUserID(uint64(1), context.Background()).
		Times(1).
		Return(retOrder, sql.ErrNoRows)

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

func TestCloseOrder(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/order/1/close", nil)

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
	handler := http.HandlerFunc(handle.CloseOrder)
	mockUseCase.EXPECT().
		CloseOrder(uint64(1), context.Background()).
		Times(1).
		Return(nil)

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

func TestCloseOrderErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/order/1/close", nil)

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
	handler := http.HandlerFunc(handle.CloseOrder)
	mockUseCase.EXPECT().
		CloseOrder(uint64(1), context.Background()).
		Times(1).
		Return(sql.ErrNoRows)

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

func TestCloseOrderErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/order/1/close", nil)

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
	handler := http.HandlerFunc(handle.CloseOrder)

	vars := map[string]string{
		"kek": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestGetAllArchiveUserOrders(t *testing.T) {
	metric.New()

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

	req, err := http.NewRequest("GET", "/api/order/profile/1/archive", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxExecutor, true)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetAllArchiveUserOrders)
	mockUseCase.EXPECT().
		GetArchiveOrders(models.UserBasicInfo{ID: 1, Executor: true}, context.Background()).
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
	metric.Destroy()
}

func TestGetAllArchiveUserOrdersErr(t *testing.T) {
	metric.New()

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

	req, err := http.NewRequest("GET", "/api/order/profile/1/archive", nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxExecutor, true)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.GetAllArchiveUserOrders)
	mockUseCase.EXPECT().
		GetArchiveOrders(models.UserBasicInfo{ID: 1, Executor: true}, context.Background()).
		Times(1).
		Return(retOrder, sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestSearchOrder(t *testing.T) {
	metric.New()

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

	keyword := models.OrderSearch{
		Keyword: "kek",
	}

	body, _ := json.Marshal(keyword)
	req, err := http.NewRequest("GET", "/api/order/profile/1/archive", bytes.NewBuffer(body))
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxExecutor, true)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SearchOrder)
	mockUseCase.EXPECT().
		SearchOrders("kek", context.Background()).
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
	metric.Destroy()
}

func TestSearchOrderErr(t *testing.T) {
	metric.New()

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

	keyword := models.OrderSearch{
		Keyword: "kek",
	}

	body, _ := json.Marshal(keyword)
	req, err := http.NewRequest("GET", "/api/order/profile/1/archive", bytes.NewBuffer(body))
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxExecutor, true)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SearchOrder)
	mockUseCase.EXPECT().
		SearchOrders("kek", context.Background()).
		Times(1).
		Return(retOrder, sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestSearchOrderJsonErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	byte22 := "meem"
	req, err := http.NewRequest("GET", "/api/order/profile/1/archive", bytes.NewBuffer([]byte(byte22)))
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()
	var val1 uint64
	var val2 uint64
	val1 = 1
	val2 = 2281488
	ctx = context.WithValue(ctx, ctxUserID, val1)
	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
	ctx = context.WithValue(ctx, ctxExecutor, true)
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SearchOrder)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}
