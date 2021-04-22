package usecase

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"post/internal/app/models"
	orderRepo "post/internal/app/order/repository"
)

const (
	orderUseCaseError = "Order use case error"
)

type UseCase struct {
	repo orderRepo.Repository
}

func NewUseCase(repo orderRepo.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Create(order models.Order) (*models.Order, error) {
	if err := u.validateOrder(&order); err != nil {
		return nil, err
	}
	u.sanitizeOrder(&order)
	var err error
	id, err := u.repo.Create(order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	order.ID = id
	err = u.supplementingTheOrderModel(&order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	return &order, err
}

func (u *UseCase) FindByID(id uint64) (*models.Order, error) {
	order, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	err = u.supplementingTheOrderModel(order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	return order, err
}

func (u *UseCase) FindByExecutorID(executorID uint64) ([]models.Order, error) {
	orders, err := u.repo.FindByExecutorID(executorID)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	for _, order := range orders {
		err = u.supplementingTheOrderModel(&order)
		if err != nil {
			return nil, errors.Wrap(err, orderUseCaseError)
		}
	}
	if orders == nil {
		return []models.Order{}, nil
	}
	return orders, nil
}

func (u *UseCase) FindByCustomerID(customerID uint64) ([]models.Order, error) {
	orders, err := u.repo.FindByCustomerID(customerID)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	for i, order := range orders {
		err = u.supplementingTheOrderModel(&order)
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

func (u *UseCase) GetActualOrders() ([]models.Order, error) {
	orders, err := u.repo.GetActualOrders()
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	for i, order := range orders {
		err = u.supplementingTheOrderModel(&order)
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

func (u *UseCase) SelectExecutor(order models.Order) error {
	//TODO: grpc-запрос за юзером
	user, err := o.store.User().FindUserByID(order.ExecutorID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	//TODO: grpc-запрос за юзером
	user.Specializes, err = o.store.User().FindSpecializesByUserID(order.ExecutorID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	if user.Executor == false {
		return errors.Wrap(err, orderUseCaseError)
	}
	if user.ID == order.CustomerID {
		return errors.Wrap(err, orderUseCaseError)
	}
	err = u.repo.UpdateExecutor(order)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (u *UseCase) DeleteExecutor(order models.Order) error {
	order.ExecutorID = 0
	err := u.repo.UpdateExecutor(order)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

//TODO: вынести в отдеьлный модуль
func (u *UseCase) validateOrder(order *models.Order) error {
	err := validation.ValidateStruct(
		order,
		validation.Field(&order.OrderName, validation.Required, validation.Length(5, 300)),
		validation.Field(&order.Description, validation.Required),
		validation.Field(&order.Category, validation.Required),
	)
	return err
}

//TODO: вынести в отдельный модуль
func (u *UseCase) sanitizeOrder(order *models.Order) {
	sanitizer := bluemonday.UGCPolicy()
	order.Category = sanitizer.Sanitize(order.Category)
	order.OrderName = sanitizer.Sanitize(order.OrderName)
	order.Description = sanitizer.Sanitize(order.Description)
}

func (u *UseCase) supplementingTheOrderModel(order *models.Order) error {
	//TODO: grpc-запрос за юзером
	u, err := o.store.User().FindUserByID(order.CustomerID)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	//FIXME: Нaхуй не нужно
	order.Login = u.Login
	//TODO: grpc-запрос за имгой
	image, err := o.mediaStore.Image().GetImage(u.Img)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	order.Img = string(image)
	return nil
}
