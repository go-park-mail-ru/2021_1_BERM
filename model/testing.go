package model

import (
	"testing"
)

func TestUser(t *testing.T) *User {
	return &User{
		Id : 1,
		Email: "Mman221@gmail.com",
		Password: "sadasds123",
		UserName: "Man",
		FirstName: "Man1",
		SecondName: "Man2",
		Executor: false,
		Specializes: nil,
	}
}

func TestOrder(t *testing.T) *Order{
	return &Order{
		Id : 1,
		OrderName: "ABC4544",
		CustomerId: 5,
		Description: "qwdqdqDDSADA",
		Specializes: []string{
			"НЕНАВИЖУ",
			"ПИСАТЬ",
			"СУКА",
			"ТЕСТЫ",
		},
	}
}

func TestSession(t *testing.T) *Session{
	return &Session{
		SessionId: "1",
		UserId: 1,
	}
}