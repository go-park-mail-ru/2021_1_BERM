package order

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

	selectArchiveOrderByID = "SELECT * FROM post.archive_orders WHERE id=$1"

	selectOrderByExecutorID = "SELECT * FROM post.orders WHERE executor_id=$1"

	selectOrderByCustomerID = "SELECT * FROM post.orders WHERE customer_id=$1"

	selectOrders = "SELECT * FROM post.orders "

	selectArchiveOrdersByExecutorID = "SELECT * FROM post.archive_orders WHERE executor_id=$1"

	selectArchiveOrdersByCustomerID = "SELECT * FROM post.archive_orders WHERE customer_id=$1"

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
                   id,
                   customer_id, 
                   executor_id, 
                   order_name, 
                   category, 
                   budget, 
                   deadline,
                   description,
                   is_archived
		)
        VALUES (
            $1, 
            $2, 
            $3,
            $4,
            $5,
			$6,
            $7,
            $8,
        	$9
                ) RETURNING id`

	searchOrdersInTitle = "SELECT * FROM post.orders WHERE to_tsvector(order_name) @@ to_tsquery($1)"

	searchOrdersInText = "SELECT * FROM post.orders WHERE to_tsvector(description) @@ to_tsquery($1)"
	getActualOrders    = "SELECT * FROM post.orders " +
		"WHERE CASE WHEN $1 != 0 THEN budget >= $1 ELSE true END " +
		"AND CASE WHEN $2 != 0  THEN budget <= $2 ELSE true END " +
		"AND CASE WHEN $3 != '~' THEN to_tsvector(order_name) @@ to_tsquery($3) ELSE true END " +
		"AND CASE WHEN $4 != '~' THEN category = $4 ELSE true END " +
		"ORDER BY budget LIMIT $5 OFFSET $6"
	getActualOrdersDesk = "SELECT * FROM post.orders " +
		"WHERE CASE WHEN $1 != 0 THEN budget >= $1 ELSE true END " +
		"AND CASE WHEN $2 != 0  THEN budget <= $2 ELSE true END " +
		"AND CASE WHEN $3 != '~' THEN to_tsvector(order_name) @@ to_tsquery($3) ELSE true END " +
		"AND CASE WHEN $4 != '~' THEN category = $4 ELSE true END " +
		"ORDER BY budget DESC LIMIT $5 OFFSET $6"

	selectTittle = `SELECT order_name FROM post.orders WHERE order_name LIKE $1 LIMIT 5`

	selectAllTittle = `SELECT order_name FROM post.orders LIMIT 5`
)
const (
	ctxQueryParams uint8 = 4
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
	tx, err := r.db.Beginx()
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	_, err = tx.NamedExec(updateOrder, &order)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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
	param := ctx.Value(ctxQueryParams).(map[string]interface{})
	category := param["category"].(string)
	limit := param["limit"].(int)
	offset := param["offset"].(int)
	desk := param["desc"].(bool)
	budgetFrom := param["from"].(int)
	budgetTo := param["to"].(int)
	searchStr := param["search_str"].(string)
	if desk {
		if err := r.db.Select(&orders, getActualOrdersDesk, budgetFrom, budgetTo, searchStr, category, limit, offset); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
	} else {
		if err := r.db.Select(&orders, getActualOrders, budgetFrom, budgetTo, searchStr, category, limit, offset); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
	}
	return orders, nil
}

//
//- search_str: fffff
//
//- budget-from: 300
//- budget-to: 400
//
//- sort: {title, salary}
//- desc: {true, false} - сортировка по возрастанию\убыванию
//- category: 'Использование человеческих ресурсов'
//- suggest: 'Использование человеческих ресурсов'
//- limit: 1
//- offset: 25

func (r *Repository) UpdateExecutor(order models.Order, ctx context.Context) error {
	tx, err := r.db.Beginx()
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	_, err = tx.NamedExec(updateExecutor, &order)
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

func (r *Repository) CreateArchive(order models.Order, ctx context.Context) error {
	order.IsArchived = true
	_, err := r.db.Query(
		insertArchiveOrder,
		order.ID,
		order.CustomerID,
		order.ExecutorID,
		order.OrderName,
		order.Category,
		order.Budget,
		order.Deadline,
		order.Description,
		order.IsArchived)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) GetArchiveOrdersByExecutorID(executorID uint64, ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectArchiveOrdersByExecutorID, executorID); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return orders, nil
}

func (r *Repository) GetArchiveOrdersByCustomerID(customerID uint64, ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Select(&orders, selectArchiveOrdersByCustomerID, customerID); err != nil {
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
	if err := r.db.Select(&orders, searchOrdersInTitle, keyword); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	if len(orders) == 0 {
		if err := r.db.Select(&orders, searchOrdersInText, keyword); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
	}
	return orders, nil
}

func (r *Repository) FindArchiveByID(id uint64, ctx context.Context) (*models.Order, error) {
	order := models.Order{}
	if err := r.db.Get(&order, selectArchiveOrderByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return &order, nil
}

func (r *Repository) SuggestOrderTitle(suggestWord string, ctx context.Context) ([]models.SuggestOrderTitle, error) {
	var suggestTittles []models.SuggestOrderTitle
	if suggestWord == "" {
		if err := r.db.Select(&suggestTittles, selectAllTittle); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
		return suggestTittles, nil
	}
	suggestWord += "%"
	if err := r.db.Select(&suggestTittles, selectTittle, suggestWord); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return suggestTittles, nil
}
