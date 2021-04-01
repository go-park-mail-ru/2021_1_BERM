package store

import "FL_2/model"

type OrderRepository interface {
	Create(order *model.Order) error
	Find(order *model.Order) error
}

