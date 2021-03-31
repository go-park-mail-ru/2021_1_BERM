package postgresstore

import (
	"fl_ru/model"
)

type OrderRepository struct {
	store *Store
}

func (o *OrderRepository) Create(order *model.Order) error {
	return nil
}

func (o *OrderRepository) Find(order *model.Order) error {
	return nil
}
