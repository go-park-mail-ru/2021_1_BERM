package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	customErr "post/internal/app/errors"
	"post/internal/app/models"
)

const (
	duplicateErrorCode = "23505"
	sqlDbSourceError   = "SQL sb source error"
)

const (
	insertOrder = `INSERT INTO post.orders (
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

	selectOrderByID = "SELECT * FROM post.orders WHERE id=$1"

	selectOrderByExecutorID = "SELECT * FROM post.orders WHERE executor_id=$1"

	selectOrderByCustomerID = "SELECT * FROM post.orders WHERE customer_id=$1"

	selectOrders = "SELECT * FROM post.orders"

	updateExecutor = `UPDATE post.orders SET 
                 executor_id =:executor_id
				 WHERE id = :id`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(order models.Order) (uint64, error) {
	var orderID uint64
	err := r.db.QueryRow(
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
				return 0, errors.Wrap(&customErr.DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	return orderID, nil
}

func (r *Repository) FindByID(id uint64) (*models.Order, error) {
	order := models.Order{}
	if err := r.db.Get(&order, selectOrderByID, id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &order, nil
}

func (r *Repository) FindByExecutorID(executorID uint64) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectOrderByExecutorID, executorID); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (r *Repository) FindByCustomerID(customerID uint64) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectOrderByCustomerID, customerID); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (r *Repository) GetActualOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectOrders); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return orders, nil
}

func (r *Repository) UpdateExecutor(order models.Order) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(updateExecutor, &order)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
