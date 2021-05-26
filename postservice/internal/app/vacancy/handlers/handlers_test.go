package vacancy

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
	"post/internal/app/vacancy/mock"
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

	vacancy := models.Vacancy{
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
	}

	retVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
		Login:       "astlok",
	}

	body, _ := json.Marshal(vacancy)

	req, err := http.NewRequest("POST", "/api/vacancy", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.CreateVacancy)
	mockUseCase.EXPECT().
		Create(vacancy, context.Background()).
		Times(1).
		Return(retVacancy, nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected, _ := json.Marshal(retVacancy)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())

	mockUseCase.EXPECT().
		Create(vacancy, context.Background()).
		Times(1).
		Return(retVacancy, sql.ErrNoRows)

	req, err = http.NewRequest("POST", "/api/vacancy", bytes.NewBuffer(body))
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
	req, err = http.NewRequest("POST", "/api/vacancy", bytes.NewBuffer([]byte(byte22)))
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

//func TestGetActualVacancy(t *testing.T) {
//	metric.New()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUseCase := mock.NewMockUseCase(ctrl)
//
//	handle := NewHandler(mockUseCase)
//
//	retVacancy := []models.Vacancy{
//		{
//			ID:          1,
//			VacancyName: "Сверстать сайт",
//			Category:    "Back",
//			CustomerID:  1,
//			Salary:      1488,
//			Description: "Pomogite sdelat API",
//			Login:       "astlok",
//		},
//	}
//
//	req, err := http.NewRequest("GET", "/api/vacancy", nil)
//
//	ctx := req.Context()
//	var val1 uint64
//	var val2 uint64
//	val1 = 1
//	val2 = 2281488
//	ctx = context.WithValue(ctx, ctxUserID, val1)
//	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
//	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
//	req = req.WithContext(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(handle.GetActualVacancies)
//	mockUseCase.EXPECT().
//		GetActualVacancies(context.Background()).
//		Times(1).
//		Return(retVacancy, nil)
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	expected, _ := json.Marshal(retVacancy)
//	expectedStr := string(expected) + "\n"
//	require.Equal(t, expectedStr, rr.Body.String())
//	metric.Destroy()
//}

//func TestGetActualVacancyErr(t *testing.T) {
//	metric.New()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUseCase := mock.NewMockUseCase(ctrl)
//
//	handle := NewHandler(mockUseCase)
//
//	retVacancy := []models.Vacancy{
//		{
//			ID:          1,
//			VacancyName: "Сверстать сайт",
//			Category:    "Back",
//			CustomerID:  1,
//			Salary:      1488,
//			Description: "Pomogite sdelat API",
//			Login:       "astlok",
//		},
//	}
//
//	req, err := http.NewRequest("GET", "/api/vacancy", nil)
//
//	ctx := req.Context()
//	var val1 uint64
//	var val2 uint64
//	val1 = 1
//	val2 = 2281488
//	ctx = context.WithValue(ctx, ctxUserID, val1)
//	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
//	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
//	req = req.WithContext(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	vars := map[string]string{
//		"id": "1",
//	}
//	req = mux.SetURLVars(req, vars)
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(handle.GetActualVacancies)
//	mockUseCase.EXPECT().
//		GetActualVacancies(context.Background()).
//		Times(1).
//		Return(retVacancy, sql.ErrNoRows)
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusInternalServerError {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusInternalServerError)
//	}
//	metric.Destroy()
//}

func TestGetVacancy(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	retVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
		Login:       "astlok",
	}

	req, err := http.NewRequest("GET", "/api/vacancy/1", nil)

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
	handler := http.HandlerFunc(handle.GetVacancy)
	mockUseCase.EXPECT().
		FindByID(uint64(1), context.Background()).
		Times(1).
		Return(retVacancy, nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(retVacancy)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())

	req1, err := http.NewRequest("GET", "/api/vacancy/1", nil)

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

func TestGetVacancyErr(t *testing.T) {

	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	retVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
		Login:       "astlok",
	}
	req, err := http.NewRequest("GET", "/api/vacancy/1", nil)

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
	handler := http.HandlerFunc(handle.GetVacancy)
	mockUseCase.EXPECT().
		FindByID(uint64(1), context.Background()).
		Times(1).
		Return(retVacancy, sql.ErrNoRows)

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

func TestChangeVacancy(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	vacancy := models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
	}
	retVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
		Login:       "astlok",
	}
	body, _ := json.Marshal(vacancy)
	req, err := http.NewRequest("GET", "/api/vacancy/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.ChangeVacancy)
	mockUseCase.EXPECT().
		ChangeVacancy(vacancy, context.Background()).
		Times(1).
		Return(*retVacancy, nil)

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(retVacancy)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())
	metric.Destroy()
}

func TestChangeVacancyBadJson(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	byte22 := "kek"
	req, err := http.NewRequest("GET", "/api/vacancy/1", bytes.NewBuffer([]byte(byte22)))

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
	handler := http.HandlerFunc(handle.ChangeVacancy)

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

func TestChangeVacancyErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	vacancy := models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
	}
	retVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
		Login:       "astlok",
	}
	body, _ := json.Marshal(vacancy)
	req, err := http.NewRequest("GET", "/api/vacancy/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.ChangeVacancy)
	mockUseCase.EXPECT().
		ChangeVacancy(vacancy, context.Background()).
		Times(1).
		Return(*retVacancy, sql.ErrNoRows)

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

func TestChangeVacancyErrParse(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	vacancy := models.Vacancy{
		ID:          1,
		VacancyName: "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Salary:      1488,
		Description: "Pomogite sdelat API",
	}
	body, _ := json.Marshal(vacancy)
	req, err := http.NewRequest("GET", "/api/vacancy/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.ChangeVacancy)

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

func TestDeleteVacancy(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/vacancy/1", nil)

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
	handler := http.HandlerFunc(handle.DeleteVacancy)
	mockUseCase.EXPECT().
		DeleteVacancy(uint64(1), context.Background()).
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

func TestDeleteVacancyErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/vacancy/1", nil)

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
	handler := http.HandlerFunc(handle.DeleteVacancy)

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

func TestDeleteVacancyErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/vacancy/1", nil)

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
	handler := http.HandlerFunc(handle.DeleteVacancy)
	mockUseCase.EXPECT().
		DeleteVacancy(uint64(1), context.Background()).
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

	vacancy := models.Vacancy{
		ID:         1,
		ExecutorID: 1,
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/vacancy/1/select", bytes.NewBuffer(body))

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
		SelectExecutor(vacancy, context.Background()).
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
	req, err := http.NewRequest("POST", "/api/vacancy/1/select", bytes.NewBuffer([]byte(byte22)))

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

	vacancy := models.Vacancy{
		ID:         1,
		ExecutorID: 1,
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/vacancy/1/select", bytes.NewBuffer(body))

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

	vacancy := models.Vacancy{
		ID:         1,
		ExecutorID: 1,
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("POST", "/api/vacancy/1/select", bytes.NewBuffer(body))

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
		SelectExecutor(vacancy, context.Background()).
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

	vacancy := models.Vacancy{
		ID: 1,
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/vacancy/1/select", bytes.NewBuffer(body))

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
		DeleteExecutor(vacancy, context.Background()).
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
	req, err := http.NewRequest("POST", "/api/vacancy/1/select", nil)

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

	vacancy := models.Vacancy{
		ID: 1,
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/vacancy/1/select", bytes.NewBuffer(body))

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
		DeleteExecutor(vacancy, context.Background()).
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

func TestGetAllVacancys(t *testing.T) {
	metric.New()

	vacancy := models.Vacancy{
		ID: 1,
	}

	retVacancy := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Сверстать сайт",
			Category:    "Back",
			CustomerID:  1,
			Salary:      1488,
			Description: "Pomogite sdelat API",
			Login:       "astlok",
		},
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/vacancy/profile/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.GetAllUserVacancies)
	mockUseCase.EXPECT().
		FindByUserID(uint64(1), context.Background()).
		Times(1).
		Return(retVacancy, nil)

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

func TestGetAllVacancysErrVar(t *testing.T) {
	metric.New()

	vacancy := models.Vacancy{
		ID: 1,
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/vacancy/profile/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.GetAllUserVacancies)

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

func TestGetAllVacancysErr(t *testing.T) {
	metric.New()

	vacancy := models.Vacancy{
		ID: 1,
	}

	retVacancy := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Сверстать сайт",
			Category:    "Back",
			CustomerID:  1,
			Salary:      1488,
			Description: "Pomogite sdelat API",
			Login:       "astlok",
		},
	}

	body, _ := json.Marshal(vacancy)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("GET", "/api/vacancy/profile/1", bytes.NewBuffer(body))

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
	handler := http.HandlerFunc(handle.GetAllUserVacancies)
	mockUseCase.EXPECT().
		FindByUserID(uint64(1), context.Background()).
		Times(1).
		Return(retVacancy, sql.ErrNoRows)

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

func TestCloseVacancy(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/vacancy/1/close", nil)

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
	handler := http.HandlerFunc(handle.CloseVacancy)
	mockUseCase.EXPECT().
		CloseVacancy(uint64(1), context.Background()).
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

func TestCloseVacancyErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/vacancy/1/close", nil)

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
	handler := http.HandlerFunc(handle.CloseVacancy)
	mockUseCase.EXPECT().
		CloseVacancy(uint64(1), context.Background()).
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

func TestCloseVacancyErrVar(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)
	req, err := http.NewRequest("DELETE", "/api/vacancy/1/close", nil)

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
	handler := http.HandlerFunc(handle.CloseVacancy)

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

//func TestGetAllArchiveUserVacancys(t *testing.T) {
//	metric.New()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUseCase := mock.NewMockUseCase(ctrl)
//
//	handle := NewHandler(mockUseCase)
//
//	retVacancy := []models.Vacancy{
//		{
//			ID:          1,
//			VacancyName: "Сверстать сайт",
//			Category:    "Back",
//			CustomerID:  1,
//			Salary:      1488,
//			Description: "Pomogite sdelat API",
//			Login:       "astlok",
//		},
//	}
//
//	req, err := http.NewRequest("GET", "/api/vacancy/profile/1/archive", nil)
//	vars := map[string]string{
//		"id": "1",
//	}
//	req = mux.SetURLVars(req, vars)
//	ctx := req.Context()
//	var val1 uint64
//	var val2 uint64
//	val1 = 1
//	val2 = 2281488
//	ctx = context.WithValue(ctx, ctxUserID, val1)
//	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
//	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
//	req = req.WithContext(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(handle.GetAllArchiveUserVacancies)
//	mockUseCase.EXPECT().
//		GetArchiveVacancies(context.Background()).
//		Times(1).
//		Return(retVacancy, nil)
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	expected, _ := json.Marshal(retVacancy)
//	expectedStr := string(expected) + "\n"
//	require.Equal(t, expectedStr, rr.Body.String())
//	metric.Destroy()
//}

//func TestGetAllArchiveUserVacancysErr(t *testing.T) {
//	metric.New()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUseCase := mock.NewMockUseCase(ctrl)
//
//	handle := NewHandler(mockUseCase)
//
//	retVacancy := []models.Vacancy{
//		{
//			ID:          1,
//			VacancyName: "Сверстать сайт",
//			Category:    "Back",
//			CustomerID:  1,
//			Salary:      1488,
//			Description: "Pomogite sdelat API",
//			Login:       "astlok",
//		},
//	}
//
//	req, err := http.NewRequest("GET", "/api/vacancy/profile/1/archive", nil)
//	vars := map[string]string{
//		"id": "1",
//	}
//	req = mux.SetURLVars(req, vars)
//	ctx := req.Context()
//	var val1 uint64
//	var val2 uint64
//	val1 = 1
//	val2 = 2281488
//	ctx = context.WithValue(ctx, ctxUserID, val1)
//	ctx = context.WithValue(ctx, ctxKeyReqID, val2)
//	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
//	req = req.WithContext(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(handle.GetAllArchiveUserVacancies)
//	mockUseCase.EXPECT().
//		GetArchiveVacancies(context.Background()).
//		Times(1).
//		Return(retVacancy, sql.ErrNoRows)
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusInternalServerError {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusInternalServerError)
//	}
//
//	metric.Destroy()
//}

func TestSearchVacancy(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	retVacancy := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Сверстать сайт",
			Category:    "Back",
			CustomerID:  1,
			Salary:      1488,
			Description: "Pomogite sdelat API",
			Login:       "astlok",
		},
	}

	keyword := models.VacancySearch{
		Keyword: "kek",
	}

	body, _ := json.Marshal(keyword)
	req, err := http.NewRequest("GET", "/api/vacancy/profile/1/archive", bytes.NewBuffer(body))
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
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SearchVacancy)
	mockUseCase.EXPECT().
		SearchVacancy("kek", context.Background()).
		Times(1).
		Return(retVacancy, nil)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(retVacancy)
	expectedStr := string(expected) + "\n"
	require.Equal(t, expectedStr, rr.Body.String())
	metric.Destroy()
}

func TestSearchVacancyErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	retVacancy := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Сверстать сайт",
			Category:    "Back",
			CustomerID:  1,
			Salary:      1488,
			Description: "Pomogite sdelat API",
			Login:       "astlok",
		},
	}

	keyword := models.VacancySearch{
		Keyword: "kek",
	}

	body, _ := json.Marshal(keyword)
	req, err := http.NewRequest("GET", "/api/vacancy/profile/1/archive", bytes.NewBuffer(body))
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
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SearchVacancy)
	mockUseCase.EXPECT().
		SearchVacancy("kek", context.Background()).
		Times(1).
		Return(retVacancy, sql.ErrNoRows)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	metric.Destroy()
}

func TestSearchVacancyJsonErr(t *testing.T) {
	metric.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	byte22 := "meem"
	req, err := http.NewRequest("GET", "/api/vacancy/profile/1/archive", bytes.NewBuffer([]byte(byte22)))
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
	ctx = context.WithValue(ctx, ctxKeyStartReqTime, time.Now())
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.SearchVacancy)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	metric.Destroy()
}
