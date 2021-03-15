package apiserver

import (
	"bytes"
	"encoding/json"
	"fl_ru/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCorrectUserInStore(t *testing.T, server *server) {
	u := model.TestUser(t)
	bCorrectUser, err := json.Marshal(u)
	assert.NoError(t, err)
	handler := server.handleSignUp()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewReader(bCorrectUser))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateCorrectOrder(t *testing.T, server *server) {
	o := model.TestOrder(t)
	bCorrectOrder, err := json.Marshal(o)
	assert.NoError(t, err)
	handler := server.handleCreateOrder()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/order", bytes.NewReader(bCorrectOrder))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
