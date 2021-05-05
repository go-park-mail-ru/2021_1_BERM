package order

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"post/internal/app/models"
	"post/internal/app/order/mock"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewHandler(mockUseCase)

	order := models.Order{
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  322,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
	}

	retOrder := models.Order{
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  322,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
		UserLogin: "astlok",
	}



	body, _:= json.Marshal(order)

	req, err := http.NewRequest("POST", "/api/order", bytes.NewBuffer(body))

	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxKeyReqID, 2281488)
	ctx = context.WithValue(ctx, ctxUserID, 1)
	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle.CreateOrder)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	mockUseCase.EXPECT().
		Create(order, ctx).
		Times(1).
		Return(retOrder, nil)

	expected := `{"id":1,"order_name":"Сверстать сайт","customer_id":322,"budget":1488,"deadline":1617004533,"description":"Pomogite sdelat API","category":"Back","user_login":"astlok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
