package order

import (
	"context"
	"fmt"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"post/internal/app/models"
	"reflect"
	"testing"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	restOrder := models.Order{
		CustomerID:  1,
		ExecutorID:  1,
		OrderName:   "Vasya",
		Category:    "Web",
		Budget:      1488,
		Deadline:    81488322,
		Description: "kekmemlul",
	}
	rows := sqlxmock.
		NewRows([]string{"orderID"}).AddRow(1)

	mock.
		ExpectQuery("INSERT INTO post.orders").
		WithArgs(restOrder.CustomerID,
			restOrder.ExecutorID,
			restOrder.OrderName,
			restOrder.Category,
			restOrder.Budget,
			restOrder.Deadline,
			restOrder.Description).
		WillReturnRows(rows)

	id, err := repo.Create(restOrder, context.Background())
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
		ExpectQuery("INSERT INTO post.orders").
		WithArgs(restOrder.CustomerID,
			restOrder.ExecutorID,
			restOrder.OrderName,
			restOrder.Category,
			restOrder.Budget,
			restOrder.Deadline,
			restOrder.Description).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Create(restOrder, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRepository_FindByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := NewRepo(db)

	restOrder := &models.Order{
		ID:          1,
		CustomerID:  1,
		ExecutorID:  1,
		OrderName:   "Vasya",
		Category:    "Web",
		Budget:      1488,
		Deadline:    81488322,
		Description: "kekmemlul",
	}
	rows := sqlxmock.
		NewRows([]string{"id", "customer_id", "executor_id", "order_name", "category", "budget", "deadline", "description"})
	rows.AddRow(1, 1, 1, "Vasya", "Web", 1488, 81488322, "kekmemlul")

	mock.
		ExpectQuery("SELECT").
		WithArgs(restOrder.ID).
		WillReturnRows(rows)

	order, err := store.FindByID(restOrder.ID, context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(order, restOrder) {
		t.Errorf("results not match, want %v, have %v", order, restOrder)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(restOrder.ID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.FindByID(restOrder.ID, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRepository_FindByExecutorID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := NewRepo(db)

	restOrder := []models.Order{
		{ID: 1,
			CustomerID:  1,
			ExecutorID:  1,
			OrderName:   "Vasya",
			Category:    "Web",
			Budget:      1488,
			Deadline:    81488322,
			Description: "kekmemlul"},
	}
	rows := sqlxmock.
		NewRows([]string{"id", "customer_id", "executor_id", "order_name", "category", "budget", "deadline", "description"})
	rows.AddRow(1, 1, 1, "Vasya", "Web", 1488, 81488322, "kekmemlul")

	mock.
		ExpectQuery("SELECT").
		WithArgs(restOrder[0].ExecutorID).
		WillReturnRows(rows)

	order, err := store.FindByExecutorID(restOrder[0].ExecutorID, context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(order, restOrder) {
		t.Errorf("results not match, want %v, have %v", order, restOrder)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(restOrder[0].ExecutorID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.FindByExecutorID(restOrder[0].ExecutorID, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRepository_FindByCustomerID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := NewRepo(db)

	restOrder := []models.Order{
		{ID: 1,
			CustomerID:  1,
			ExecutorID:  1,
			OrderName:   "Vasya",
			Category:    "Web",
			Budget:      1488,
			Deadline:    81488322,
			Description: "kekmemlul"},
	}
	rows := sqlxmock.
		NewRows([]string{"id", "customer_id", "executor_id", "order_name", "category", "budget", "deadline", "description"})
	rows.AddRow(1, 1, 1, "Vasya", "Web", 1488, 81488322, "kekmemlul")

	mock.
		ExpectQuery("SELECT").
		WithArgs(restOrder[0].CustomerID).
		WillReturnRows(rows)

	order, err := store.FindByCustomerID(restOrder[0].CustomerID, context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(order, restOrder) {
		t.Errorf("results not match, want %v, have %v", order, restOrder)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(restOrder[0].CustomerID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.FindByCustomerID(restOrder[0].CustomerID, context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestRepository_GetActualOrders(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	store := NewRepo(db)

	restOrder := []models.Order{
		{ID: 1,
			CustomerID:  1,
			ExecutorID:  1,
			OrderName:   "Vasya",
			Category:    "Web",
			Budget:      1488,
			Deadline:    81488322,
			Description: "kekmemlul"},
	}
	rows := sqlxmock.
		NewRows([]string{"id", "customer_id", "executor_id", "order_name", "category", "budget", "deadline", "description"})
	rows.AddRow(1, 1, 1, "Vasya", "Web", 1488, 81488322, "kekmemlul")

	mock.
		ExpectQuery("SELECT").
		WillReturnRows(rows)

	order, err := store.GetActualOrders(context.Background())

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(order, restOrder) {
		t.Errorf("results not match, want %v, have %v", order, restOrder)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.
		ExpectQuery("SELECT").
		WillReturnError(fmt.Errorf("db_error"))

	_, err = store.GetActualOrders(context.Background())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
