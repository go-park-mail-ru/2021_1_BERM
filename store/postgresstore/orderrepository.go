package postgresstore

import (
	"FL_2/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type OrderRepository struct {
	store *Store
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

func (o *OrderRepository) Create(order model.Order) (uint64, error) {
	var orderID uint64
	err := o.store.db.QueryRow(
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
			if pqErr.Code == duplicateErrorCode {
				return 0, errors.Wrap(&DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	return orderID, nil
}

func (o *OrderRepository) FindByID(id uint64) (*model.Order, error) {
	order := model.Order{}
	if err := o.store.db.Get(&order, selectOrderByID, id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &order, nil
}

func (o *OrderRepository) FindByExecutorID(executorID uint64) ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, selectOrderByExecutorID, executorID); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) FindByCustomerID(customerID uint64) ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, selectOrderByCustomerID, customerID); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) GetActualOrders() ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, selectOrders); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) UpdateExecutor(order model.Order) error {
	tx := o.store.db.MustBegin()
	_, err := tx.NamedExec(updateExecutor, &order)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
