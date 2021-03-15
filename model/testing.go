package model

import (
	"testing"
)

const (
	CustomerID uint64 = 5
)

func TestUser(t *testing.T) *User {
	return &User{
		ID:          1,
		Email:       "Mman221@gmail.com",
		Password:    "sadasds123",
		UserName:    "Man",
		FirstName:   "Man1",
		SecondName:  "Man2",
		Executor:    false,
		Specializes: nil,
	}
}

func TestOrder(t *testing.T) *Order {
	return &Order{
		ID:          1,
		OrderName:   "ABC4544",
		CustomerID:  CustomerID,
		Description: "qwdqdqDDSADA",
		Specializes: []string{
			"НЕНАВИЖУ",
			"ПИСАТЬ",
			"СУКА",
			"ТЕСТЫ",
		},
	}
}

func TestSession(t *testing.T) *Session {
	return &Session{
		SessionID: "1",
		UserID:    1,
	}
}
