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
