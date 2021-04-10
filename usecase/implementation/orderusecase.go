package implementation

import (
	"FL_2/model"
	"FL_2/store"
	validation "github.com/go-ozzo/ozzo-validation"
)

type OrderUseCase struct {
	store store.Store
}

func (o *OrderUseCase) Create(order model.Order) (uint64, error){
	if err := o.validateOrder(&order); err != nil {
		return 0, err
	}
	var err error
	id, err := o.store.Order().Create(order)
	return id, err
}

func (o *OrderUseCase) FindByID(id uint64) (*model.Order, error){
	order, err := o.store.Order().FindByID(id)
	return order, err
}

func (o *OrderUseCase) FindByExecutorID(executorID uint64) ([]model.Order, error){
	order, err := o.store.Order().FindByExecutorID(executorID)
	return order, err
}

func (o *OrderUseCase) FindByCustomerID(customerID uint64) ([]model.Order, error){
	order, err := o.store.Order().FindByCustomerID(customerID)
	return order, err
}

func (o *OrderUseCase) GetActualOrders()([]model.Order, error){
	orders, err := o.store.Order().GetActualOrders()
	return orders, err
}

func (o *OrderUseCase)validateOrder(order *model.Order) error{
	err := validation.ValidateStruct(
		order,
		validation.Field(&order.OrderName, validation.Required,  validation.Length(5, 300)),
		validation.Field(&order.Description, validation.Required,),
		validation.Field(&order.Category, validation.Required),
	)
	return err;
}
