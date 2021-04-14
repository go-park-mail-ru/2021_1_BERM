package implementation

import (
	"FL_2/model"
	"FL_2/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
)

const (
	orderUseCaseError = "Order use case error"
)

type OrderUseCase struct {
	store      store.Store
	mediaStore store.MediaStore
}

func (o *OrderUseCase) Create(order model.Order) (*model.Order, error) {
	if err := o.validateOrder(&order); err != nil {
		return nil, err
	}
	o.sanitizeOrder(&order)
	var err error
	id, err := o.store.Order().Create(order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	order.ID = id
	err = o.supplementingTheOrderModel(&order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	return &order, err
}

func (o *OrderUseCase) FindByID(id uint64) (*model.Order, error) {
	order, err := o.store.Order().FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	err = o.supplementingTheOrderModel(order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	return order, err
}

func (o *OrderUseCase) FindByExecutorID(executorID uint64) ([]model.Order, error) {
	orders, err := o.store.Order().FindByExecutorID(executorID)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	for _, order := range orders {
		err = o.supplementingTheOrderModel(&order)
		if err != nil {
			return nil, errors.Wrap(err, orderUseCaseError)
		}
	}
	if orders == nil {
		return []model.Order{}, nil
	}
	return orders, nil
}

func (o *OrderUseCase) FindByCustomerID(customerID uint64) ([]model.Order, error) {
	orders, err := o.store.Order().FindByCustomerID(customerID)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	for i, order := range orders {
		err = o.supplementingTheOrderModel(&order)
		if err != nil {
			return nil, errors.Wrap(err, orderUseCaseError)
		}
		orders[i] = order
	}
	if orders == nil {
		return []model.Order{}, nil
	}
	return orders, nil
}

func (o *OrderUseCase) GetActualOrders() ([]model.Order, error) {
	orders, err := o.store.Order().GetActualOrders()
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	for i, order := range orders {
		err = o.supplementingTheOrderModel(&order)
		if err != nil {
			return nil, errors.Wrap(err, orderUseCaseError)
		}
		orders[i] = order
	}
	if orders == nil {
		return []model.Order{}, nil
	}
	return orders, err
}

func (o *OrderUseCase) SelectExecutor(order model.Order) error {
	user, err := o.store.User().FindByID(order.ExecutorID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	if user.Executor == false {
		return errors.Wrap(err, orderUseCaseError)
	}
	if user.ID == order.CustomerID {
		return errors.Wrap(err, orderUseCaseError)
	}
	err = o.store.Order().AddExecutor(order)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (o *OrderUseCase) validateOrder(order *model.Order) error {
	err := validation.ValidateStruct(
		order,
		validation.Field(&order.OrderName, validation.Required, validation.Length(5, 300)),
		validation.Field(&order.Description, validation.Required),
		validation.Field(&order.Category, validation.Required),
	)
	return err
}

func (o *OrderUseCase) supplementingTheOrderModel(order *model.Order) error {
	u, err := o.store.User().FindByID(order.CustomerID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	order.Login = u.Login
	image, err := o.mediaStore.Image().GetImage(u.Img)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	order.Img = string(image)
	return nil
}

func (o *OrderUseCase) sanitizeOrder(order *model.Order) {
	sanitizer := bluemonday.UGCPolicy()
	order.Category = sanitizer.Sanitize(order.Category)
	order.OrderName = sanitizer.Sanitize(order.OrderName)
	order.Description = sanitizer.Sanitize(order.Description)
}
