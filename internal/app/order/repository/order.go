package repository

import (
	"ff/internal/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	ffError "ff/internal/app/server/errors"
)

type OrderRepository struct {
	db *sqlx.DB
}

const (
	insertOrder = `INSERT INTO ff.orders (
                   customer_id, 
                   executor_id, 
                   order_name, 
                   category, 
                   budget, 
                   deadline,
                   description
		)
        VALUES (
                $1, 
                $2, 
                $3,
				$4,
				$5,
				$6,
                $7
                ) RETURNING id`

	selectOrderByID = "SELECT * FROM ff.orders WHERE id=$1"

	selectOrderByExecutorID = "SELECT * FROM ff.orders WHERE executor_id=$1"

	selectOrderByCustomerID = "SELECT * FROM ff.orders WHERE customer_id=$1"

	selectOrders = "SELECT * FROM ff.orders"

	updateExecutor = `UPDATE ff.orders SET 
                 executor_id =:executor_id
				 WHERE id = :id`
)

func (o *OrderRepository) Create(order models.Order) (uint64, error) {
	var orderID uint64
	err := o.db.QueryRow(
		insertOrder,
		order.CustomerID,
		order.ExecutorID,
		order.OrderName,
		order.Category,
		order.Budget,
		order.Deadline,
		order.Description).Scan(&orderID)
	if err != nil {
		pqErr := &pq.Error{}
		if errors.As(err, &pqErr) {
			if pqErr.Code == ffError.DuplicateErrorCode {
				return 0, errors.Wrap(&ffError.DuplicateSourceErr{
					Err: err,
				}, ffError.SqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, ffError.SqlDbSourceError)
	}
	return orderID, nil
}

func (o *OrderRepository) FindByID(id uint64) (*models.Order, error) {
	order := models.Order{}
	if err := o.db.Get(&order, selectOrderByID, id); err != nil {
		return nil, errors.Wrap(err, ffError.SqlDbSourceError)
	}
	return &order, nil
}

func (o *OrderRepository) FindByExecutorID(executorID uint64) ([]models.Order, error) {
	var orders []models.Order
	if err := o.db.Select(&orders, selectOrderByExecutorID, executorID); err != nil {
		return nil, errors.Wrap(err, ffError.SqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) FindByCustomerID(customerID uint64) ([]models.Order, error) {
	var orders []models.Order
	if err := o.db.Select(&orders, selectOrderByCustomerID, customerID); err != nil {
		return nil, errors.Wrap(err, ffError.SqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) GetActualOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := o.db.Select(&orders, selectOrders); err != nil {
		return nil, errors.Wrap(err, ffError.SqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) UpdateExecutor(order models.Order) error {
	tx := o.db.MustBegin()
	_, err := tx.NamedExec(updateExecutor, &order)
	if err != nil {
		return errors.Wrap(err, ffError.SqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, ffError.SqlDbSourceError)
	}
	return nil
}
