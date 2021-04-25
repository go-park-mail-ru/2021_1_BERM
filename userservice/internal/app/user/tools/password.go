package tools

import (
	"bytes"
	"golang.org/x/crypto/argon2"
	"math/rand"
	"user/Error"
	"user/internal/app/models"
)

const(
	saltLength = 8
)

func CompPass(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[0:8])

	userPassHash := hashPass(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}

func BeforeCreate(user *models.NewUser) error {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return &Error.Error{
			Err: err,
			InternalError: true,
			ErrorDescription: Error.InternalServerErrorDescription,
		}
	}

	user.EncryptPassword = hashPass(salt, user.Password)
	if user.Specializes != nil{
		user.Executor = true
	}
	return nil
}

func BeforeChange(user *models.ChangeUser) error {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return &Error.Error{
			Err: err,
			InternalError: true,
			ErrorDescription: Error.InternalServerErrorDescription,
		}
	}

	user.EncryptPassword = hashPass(salt, user.Password)
	if user.Specializes != nil{
		user.Executor = true
	}
	return nil
}


func hashPass(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}



