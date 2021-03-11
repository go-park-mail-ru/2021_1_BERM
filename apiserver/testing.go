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

func TestCreateCorrectUserInStore(t *testing.T, server *server){
	u := model.TestUser(t)
	bCorrectUser, err := json.Marshal(u)
	assert.NoError(t, err)
	handler := http.HandlerFunc(server.handleSignUp())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewReader(bCorrectUser))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
