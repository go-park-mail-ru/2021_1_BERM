package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"post/internal/app/models"
	"post/pkg/error/errortools"
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

	selectArchiveOrders = "SELECT * FROM post.archive_orders"

	updateExecutor = `UPDATE post.orders SET 
                 executor_id =:executor_id
				 WHERE id = :id`

	updateOrder = `UPDATE post.orders SET
					order_name =:order_name,
					category =:category,
					customer_id =:customer_id,
					executor_id =:executor_id,
					deadline =:deadline,
					budget =:budget,
					description =:description
					WHERE id =:id`

	deleteOrder = `DELETE from post.orders WHERE id=$1`

	insertArchiveOrder = `INSERT INTO post.archive_orders (
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

	searchOrders = "SELECT * FROM post.orders WHERE to_tsvector(order_name) @@ to_tsquery($1)"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(order models.Order, ctx context.Context) (uint64, error) {
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
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return orderID, nil
}

func (r *Repository) Change(order models.Order, ctx context.Context) error {
	tx, err := r.db.Begin()
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	_, err = tx.Exec(updateOrder, order)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	if err = tx.Commit(); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) DeleteOrder(id uint64, ctx context.Context) error {
	_, err := r.db.Queryx(deleteOrder, id)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) FindByID(id uint64, ctx context.Context) (*models.Order, error) {
	order := models.Order{}
	if err := r.db.Get(&order, selectOrderByID, id); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return &order, nil
}

func (r *Repository) FindByExecutorID(executorID uint64, ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Select(&orders, selectOrderByExecutorID, executorID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return orders, nil
}

func (r *Repository) FindByCustomerID(customerID uint64, ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Select(&orders, selectOrderByCustomerID, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return orders, nil
}

func (r *Repository) GetActualOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectOrders); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return orders, nil
}

func (r *Repository) UpdateExecutor(order models.Order, ctx context.Context) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(updateExecutor, &order)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	if err := tx.Commit(); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) CreateArchive(order models.Order, ctx context.Context) (uint64, error) {
	var orderID uint64
	err := r.db.QueryRow(
		insertArchiveOrder,
		order.CustomerID,
		order.ExecutorID,
		order.OrderName,
		order.Category,
		order.Budget,
		order.Deadline,
		order.Description).Scan(&orderID)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return orderID, nil
}

func (r *Repository) GetArchiveOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectArchiveOrders); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return orders, nil
}

func (r *Repository) SearchOrders(keyword string, ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	if keyword == "" {
		return nil, nil
	}
	keyword += ":*"
	if err := r.db.Select(&orders, searchOrders, keyword); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return orders, nil
}
