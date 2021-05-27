package repository_test

import (
	"context"
	"fmt"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"reflect"
	"testing"
	"user/internal/app/models"
	specializeRep "user/internal/app/specialize/repository"
)

func TestFindSpecializeByName(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := specializeRep.Repository{
		Db: db,
	}
	testSpec := models.Specialize{
		ID:   1,
		Name: "Back",
	}

	testID := uint64(1)
	rows := sqlxmock.
		NewRows([]string{"id", "specialize_name"})
	rows.AddRow(1, "Back")

	mock.
		ExpectQuery("SELECT").
		WithArgs(testSpec.Name).
		WillReturnRows(rows)

	specID, err := store.FindByName(testSpec.Name, context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(specID, testID) {
		t.Errorf("results not match, want %v, have %v", specID, testSpec)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(testSpec.Name).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.FindByName(testSpec.Name, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
