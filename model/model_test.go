package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := TestUser(t)
	assert.NoError(t, u.BeforeCreate())
}

func TestUser_Validate	(t *testing.T) {
	u := TestUser(t)
	assert.NoError(t, u.Validate())
}


func TestUser_EncryptString(t *testing.T) {
	encrypt, err := encryptString("122131", "sdsadasd")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypt)
}

func TestUser_Sanitize(t *testing.T) {
	u := TestUser(t)
	u.Sanitize()
	assert.Empty(t, u.Password)
}

func TestUser_ComparePassword(t *testing.T) {
	u := TestUser(t)
	noEncryptPass := u.Password
	_ = u.BeforeCreate()
	assert.Equal(t, u.ComparePassword(noEncryptPass), true)
}

func TestOrder_Validate(t *testing.T) {
	o := TestOrder(t)
	assert.NoError(t, o.Validate())
}

func TestSession_BeforeChange(t *testing.T) {
	s := TestSession(t)
	noEncryptSession := s.SessionId
	s.BeforeChange()
	assert.NotEqual(t, noEncryptSession, s.SessionId)
}
