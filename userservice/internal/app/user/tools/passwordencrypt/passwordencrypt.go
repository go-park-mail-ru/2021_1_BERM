package passwordencrypt

import (
	"bytes"
	"golang.org/x/crypto/argon2"
	"math/rand"
	"user/internal/app/models"
)

const (
	saltLength = 8
)

type PasswordEncrypter struct {
}

func (p PasswordEncrypter) CompPass(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[0:8])

	userPassHash := p.hashPass(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}

func (p PasswordEncrypter) BeforeCreate(user models.NewUser) (models.NewUser, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return models.NewUser{}, err
	}

	user.EncryptPassword = p.hashPass(salt, user.Password)
	if user.Specializes != nil {
		user.Executor = true
	}
	return user, nil
}

func (p PasswordEncrypter) BeforeChange(user models.ChangeUser) (models.ChangeUser, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return models.ChangeUser{}, err
	}

	user.EncryptPassword = p.hashPass(salt, user.Password)
	if user.Specializes != nil {
		user.Executor = true
	}
	return user, nil
}

func (p PasswordEncrypter) hashPass(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}
