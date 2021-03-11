package teststore

import "fl_ru/model"

type OrderRepository struct {
	store *Store
	order map[uint64]*model.Order
}

func (r *OrderRepository)Create(order *model.Order) error{
	order.Id = 1
	r.order[1] = order

	return nil
}

func (r *OrderRepository)Find(order *model.Order) error{

	return nil
}
