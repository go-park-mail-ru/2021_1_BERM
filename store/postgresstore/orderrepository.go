package postgresstore

import (
	"FL_2/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type OrderRepository struct {
	store *Store
}

func (o *OrderRepository) Create(order model.Order) (uint64, error) {
	var orderID uint64
	err := o.store.db.QueryRow(
		`INSERT INTO orders (
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
                ) RETURNING id`,
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
	if err := o.store.db.Get(&order, "SELECT * FROM orders WHERE id=$1", id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &order, nil
}

func (o *OrderRepository) FindByExecutorID(executorID uint64) ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, "SELECT * FROM orders WHERE executor_id=$1", executorID); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) FindByCustomerID(customerID uint64) ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, "SELECT * FROM orders WHERE customer_id=$1", customerID); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) GetActualOrders() ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, "SELECT * FROM orders"); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (o *OrderRepository) UpdateExecutor(order model.Order) error {
	tx := o.store.db.MustBegin()
	_, err := tx.NamedExec(`UPDATE orders SET 
                 executor_id =:executor_id
				 WHERE id = :id`, &order)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
