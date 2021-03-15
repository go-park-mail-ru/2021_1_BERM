package tarantoolstore

import "fl_ru/model"

type OrderRepository struct {
	store *Store
}

func (o *OrderRepository) Create(order *model.Order) error {
	resp, err := o.store.conn.Insert("order", orderToTarantoolData(order))
	if err != nil {
		return err
	}
	*order = *tarantoolDataToOrder(resp.Tuples()[0])

	return nil
}

func (o *OrderRepository) Find(order *model.Order) error {
	return nil
}

func orderToTarantoolData(order *model.Order) []interface{} {
	return []interface{}{
		nil,
		order.OrderName,
		order.CustomerID,
		order.Description,
		order.Specializes,
	}
}

func tarantoolDataToOrder(data []interface{}) *model.Order {
	order := &model.Order{}
	order.ID, _ = data[0].(uint64)
	order.OrderName, _ = data[1].(string)
	order.CustomerID, _ = data[2].(uint64)
	order.Description, _ = data[3].(string)
	order.Specializes = []string{}
	specializes, _ := data[4].([]interface{})
	for _, elem := range specializes {
		specialize, _ := elem.(string)
		order.Specializes = append(order.Specializes, specialize)
	}

	return order
}
