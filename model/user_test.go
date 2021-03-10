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