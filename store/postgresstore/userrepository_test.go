package postgresstore

import (
	"FL_2/model"
	"fmt"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"reflect"
	"testing"
)

func TestUserCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := &Store{
		Db: db,
	}

	testUser := model.User{
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
		ExpectQuery("INSERT INTO ff.users").
		WithArgs(testUser.Email,
			testUser.EncryptPassword,
			testUser.Login,
			testUser.NameSurname,
			testUser.About,
			testUser.Executor).
		WillReturnRows(rows)

	id, err := store.User().AddUser(&testUser)
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
		ExpectQuery("INSERT INTO ff.users").
		WithArgs(testUser.Email,
			testUser.EncryptPassword,
			testUser.Login,
			testUser.NameSurname,
			testUser.About,
			testUser.Executor).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.User().AddUser(&testUser)
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

	store := &Store{
		Db: db,
	}

	testUser := &model.User{
		ID:              1,
		Email:           "kek@me.ru",
		EncryptPassword: []byte("123"),
		Login:           "Vasya",
		NameSurname:     "Vasya Pupkin",
		About:           "Ya Vasya",
		Executor:        true,
	}
	rows := sqlxmock.
		NewRows([]string{"id", "email", "password", "login", "name_surname", "about", "executor"})
	rows.AddRow(1, "kek@me.ru", []byte("123"), "Vasya", "Vasya Pupkin", "Ya Vasya", true)

	mock.
		ExpectQuery("SELECT").
		WithArgs(testUser.ID).
		WillReturnRows(rows)

	user, err := store.User().FindUserByID(testUser.ID)

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

	_, err = store.User().FindUserByID(testUser.ID)
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

	store := &Store{
		Db: db,
	}

	testUser := &model.User{
		ID:              1,
		Email:           "kek@me.ru",
		EncryptPassword: []byte("123"),
		Login:           "Vasya",
		NameSurname:     "Vasya Pupkin",
		About:           "Ya Vasya",
		Executor:        true,
	}
	rows := sqlxmock.
		NewRows([]string{"id", "email", "password", "login", "name_surname", "about", "executor"})
	rows.AddRow(1, "kek@me.ru", []byte("123"), "Vasya", "Vasya Pupkin", "Ya Vasya", true)

	mock.
		ExpectQuery("SELECT").
		WithArgs(testUser.Email).
		WillReturnRows(rows)

	user, err := store.User().FindUserByEmail(testUser.Email)

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

	_, err = store.User().FindUserByEmail(testUser.Email)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestFindSpecializeByName(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := &Store{
		Db: db,
	}

	testSpec := model.Specialize{
		ID:   1,
		Name: "Back",
	}
	rows := sqlxmock.
		NewRows([]string{"id", "specialize_name"})
	rows.AddRow(1, "Back")

	mock.
		ExpectQuery("SELECT").
		WithArgs(testSpec.Name).
		WillReturnRows(rows)

	spec, err := store.User().FindSpecializeByName(testSpec.Name)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(spec, testSpec) {
		t.Errorf("results not match, want %v, have %v", spec, testSpec)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(testSpec.Name).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.User().FindSpecializeByName(testSpec.Name)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
