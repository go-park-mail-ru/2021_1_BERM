package usecase

import (
	"ff/internal/app/image"
	"ff/internal/app/models"
	"ff/internal/app/order"
	user "ff/internal/app/user"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
)

const (
	orderUseCaseError = "Order use case errors"
)

type OrderUseCase struct {
	orderRepo      order.OrderRepository
	mediaStore     image.ImageRepository
}

func (o *OrderUseCase) Create(order models.Order) (*models.Order, error) {
	if err := o.validateOrder(&order); err != nil {
		return nil, err
	}
	o.sanitizeOrder(&order)
	var err error
	id, err := o.orderRepo.Create(order)
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

func (o *OrderUseCase) FindByID(id uint64) (*models.Order, error) {
	order, err := o.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	err = o.supplementingTheOrderModel(order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	return order, err
}

func (o *OrderUseCase) FindByExecutorID(executorID uint64) ([]models.Order, error) {
	orders, err := o.orderRepo.FindByExecutorID(executorID)
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
		return []models.Order{}, nil
	}
	return orders, nil
}

func (o *OrderUseCase) FindByCustomerID(customerID uint64) ([]models.Order, error) {
	orders, err := o.orderRepo.FindByCustomerID(customerID)
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
		return []models.Order{}, nil
	}
	return orders, nil
}

func (o *OrderUseCase) GetActualOrders() ([]models.Order, error) {
	orders, err := o.orderRepo.GetActualOrders()
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
		return []models.Order{}, nil
	}
	return orders, err
}

func (o *OrderUseCase) SelectExecutor(order models.Order) error {
	u, err := user.UserUseCase.FindByID(order.ExecutorID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	u.Specializes, err = user.UserRepository.FindSpecializesByUserID(order.ExecutorID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	if u.Executor == false {
		return errors.Wrap(err, orderUseCaseError)
	}
	if u.ID == order.CustomerID {
		return errors.Wrap(err, orderUseCaseError)
	}
	err = o.orderRepo.UpdateExecutor(order)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (o *OrderUseCase) DeleteExecutor(order models.Order) error {
	order.ExecutorID = 0
	err := o.orderRepo.UpdateExecutor(order)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (o *OrderUseCase) validateOrder(order *models.Order) error {
	err := validation.ValidateStruct(
		order,
		validation.Field(&order.OrderName, validation.Required, validation.Length(5, 300)),
		validation.Field(&order.Description, validation.Required),
		validation.Field(&order.Category, validation.Required),
	)
	return err
}

func (o *OrderUseCase) supplementingTheOrderModel(order *models.Order) error {
	u, err := user.UserRepository.FindUserByID(order.CustomerID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	order.Login = u.Login
	i, err := o.mediaStore.GetImage(u.Img)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	order.Img = string(i)
	return nil
}

func (o *OrderUseCase) sanitizeOrder(order *models.Order) {
	sanitizer := bluemonday.UGCPolicy()
	order.Category = sanitizer.Sanitize(order.Category)
	order.OrderName = sanitizer.Sanitize(order.OrderName)
	order.Description = sanitizer.Sanitize(order.Description)
}

