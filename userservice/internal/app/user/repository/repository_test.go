package user

import (
	"context"
	"fmt"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"reflect"
	"testing"
	"user/internal/app/models"
)

func TestUserCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := Repository{Db: db}

	testUser := models.NewUser{
		Email:           "kek@me.ru",
		EncryptPassword: []byte("123"),
		Login:           "Vasya",
		NameSurname:     "Vasya Pupkin",
		About:           "Ya Vasya",
		Executor:        true,
	}
	rows := sqlxmock.
		NewRows([]string{"orderID"}).AddRow(1)

	mock.
		ExpectQuery("INSERT INTO userservice.users").
		WithArgs(testUser.Email,
			testUser.EncryptPassword,
			testUser.Login,
			testUser.NameSurname,
			testUser.About,
			testUser.Executor).
		WillReturnRows(rows)

	id, err := store.Create(&testUser, context.Background())
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if id != 1 {
		t.Errorf("bad id: want %v, have %v", id, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("INSERT INTO userservice.users").
		WithArgs(testUser.Email,
			testUser.EncryptPassword,
			testUser.Login,
			testUser.NameSurname,
			testUser.About,
			testUser.Executor).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.Create(&testUser, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestFindUserByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := Repository{Db: db}

	testUser := &models.UserInfo{
		ID:          1,
		Email:       "kek@me.ru",
		Password:    []byte("123"),
		Login:       "Vasya",
		NameSurname: "Vasya Pupkin",
		About:       "Ya Vasya",
		Executor:    true,
	}
	rows := sqlxmock.
		NewRows([]string{"id", "email", "password", "login", "name_surname", "about", "executor"})
	rows.AddRow(1, "kek@me.ru", []byte("123"), "Vasya", "Vasya Pupkin", "Ya Vasya", true)

	mock.
		ExpectQuery("SELECT").
		WithArgs(testUser.ID).
		WillReturnRows(rows)

	user, err := store.FindUserByID(testUser.ID, context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(user, testUser) {
		t.Errorf("results not match, want %v, have %v", user, testUser)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(testUser.ID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.FindUserByID(testUser.ID, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestFindUserByEmail(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := Repository{Db: db}

	testUser := &models.UserInfo{
		ID:          1,
		Email:       "kek@me.ru",
		Password:    []byte("123"),
		Login:       "Vasya",
		NameSurname: "Vasya Pupkin",
		About:       "Ya Vasya",
		Executor:    true,
	}
	rows := sqlxmock.
		NewRows([]string{"id", "email", "password", "login", "name_surname", "about", "executor"})
	rows.AddRow(1, "kek@me.ru", []byte("123"), "Vasya", "Vasya Pupkin", "Ya Vasya", true)

	mock.
		ExpectQuery("SELECT").
		WithArgs(testUser.Email).
		WillReturnRows(rows)

	user, err := store.FindUserByEmail(testUser.Email, context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(user, testUser) {
		t.Errorf("results not match, want %v, have %v", user, testUser)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(testUser.Email).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.FindUserByEmail(testUser.Email, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
