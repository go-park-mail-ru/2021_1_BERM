package postgresstore

import "FL_2/model"

type OrderRepository struct {
	store *Store
}

func (o *OrderRepository) Create(order *model.Order) error {
	return nil
}

func (o *OrderRepository) Find(order *model.Order) error {
	return nil
}

