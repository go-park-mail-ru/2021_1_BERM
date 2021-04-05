package postgresstore

import "FL_2/model"

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
		order.Description).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (o *OrderRepository) FindByID(id uint64) (*model.Order, error) {
	order := &model.Order{}
	if err := o.store.db.Get(&order, "SELECT * FROM orders WHERE id=$1", id); err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderRepository) FindByExecutorID(executorID uint64) ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, "SELECT * FROM orders WHERE executor_id=$1", executorID); err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderRepository) FindByCustomerID(customerID uint64) ([]model.Order, error) {
	var orders []model.Order
	if err := o.store.db.Select(&orders, "SELECT * FROM orders WHERE customer_id=$1", customerID); err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderRepository)GetActualOrders()([]model.Order, error){
	var orders []model.Order
	if err := o.store.db.Select(&orders, "SELECT * FROM orders"); err != nil {
		return nil, err
	}
	return orders, nil
}