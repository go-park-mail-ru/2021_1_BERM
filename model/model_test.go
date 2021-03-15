package model_test

import (
	"fl_ru/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
}

func TestUser_Validate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.Validate())
}

func TestUser_EncryptString(t *testing.T) {
	encrypt, err := model.EncryptString("122131", "sdsadasd")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypt)
}

func TestUser_Sanitize(t *testing.T) {
	u := model.TestUser(t)
	u.Sanitize()
	assert.Empty(t, u.Password)
}

func TestUser_ComparePassword(t *testing.T) {
	u := model.TestUser(t)
	noEncryptPass := u.Password
	_ = u.BeforeCreate()
	assert.Equal(t, u.ComparePassword(noEncryptPass), true)
}

func TestOrder_Validate(t *testing.T) {
	o := model.TestOrder(t)
	assert.NoError(t, o.Validate())
}

func TestSession_BeforeChange(t *testing.T) {
	s := model.TestSession(t)
	noEncryptSession := s.SessionID
	s.BeforeChange()
	assert.NotEqual(t, noEncryptSession, s.SessionID)
}
