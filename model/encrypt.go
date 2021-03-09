package model

import "golang.org/x/crypto/bcrypt"

func encryptString(password string, salt string) (string, error){
	b, err := bcrypt.GenerateFromPassword([]byte(password + salt), bcrypt.MinCost)
	if err != nil{
		return "", err
	}
	return string(b), nil
}
