package usecase

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"post/api"
	"post/internal/app/models"
	orderRepo "post/internal/app/order/repository"
	customErr "post/pkg/error"
)

const (
	orderUseCaseError = "Order use case error"
)

type UseCase struct {
	OrderRepo orderRepo.Repository
	UserRepo  api.UserClient
}

func NewUseCase(orderRepo orderRepo.Repository, userRepo api.UserClient) *UseCase {
	return &UseCase{
		OrderRepo: orderRepo,
		UserRepo:  userRepo,
	}
}

func (u *UseCase) Create(order models.Order, ctx context.Context) (*models.Order, error) {
	if err := u.validateOrder(&order); err != nil {
		return nil, err
	}
	u.sanitizeOrder(&order)
	var err error
	id, err := u.OrderRepo.Create(order, ctx)
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

func (u *UseCase) FindByID(id uint64, ctx context.Context) (*models.Order, error) {
	order, err := u.OrderRepo.FindByID(id, ctx)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	err = u.supplementingTheOrderModel(order)
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	return order, err
}

func (u *UseCase) FindByUserID(userID uint64, ctx context.Context) ([]models.Order, error) {
	userR, err := u.UserRepo.GetUserById(ctx, &api.UserRequest{Id: userID})
	if err != nil {
		return nil, errors.Wrap(err, orderUseCaseError)
	}
	isExecutor := userR.GetExecutor()
	var orders []models.Order
	if isExecutor {
		orders, err = u.OrderRepo.FindByExecutorID(userID, ctx)
	} else {
		orders, err = u.OrderRepo.FindByCustomerID(userID, ctx)
	}
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

func (u *UseCase) ChangeOrder(order models.Order, ctx context.Context) (models.Order, error) {
	oldOrder, err := u.OrderRepo.FindByID(order.ID, context.Background())
	if err != nil {
		return models.Order{}, errors.Wrap(err, orderUseCaseError)
	}
	if order.OrderName == "" {
		order.OrderName = oldOrder.OrderName
	}
	if order.Category == "" {
		order.Category = oldOrder.Category
	}
	if order.Description == "" {
		order.Description = oldOrder.Description
	}
	if order.Budget == 0 {
		order.Budget = oldOrder.Budget
	}
	if order.Deadline == 0 {
		order.Deadline = oldOrder.Deadline
	}
	order.CustomerID = oldOrder.CustomerID
	order.ExecutorID = oldOrder.ExecutorID
	err = u.OrderRepo.Change(order, ctx)
	if err != nil {
		return models.Order{}, errors.Wrap(err, orderUseCaseError)
	}
	err = u.supplementingTheOrderModel(&order)
	if err != nil {
		return models.Order{}, errors.Wrap(err, orderUseCaseError)
	}
	return order, nil
}

func (u *UseCase) DeleteOrder(id uint64, ctx context.Context) error {
	err := u.OrderRepo.DeleteOrder(id, ctx)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (u *UseCase) GetActualOrders(ctx context.Context) ([]models.Order, error) {
	orders, err := u.OrderRepo.GetActualOrders(ctx)
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

func (u *UseCase) SelectExecutor(order models.Order, ctx context.Context) error {
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: order.ExecutorID})
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	if userR.GetExecutor() == false {
		return customErr.ErrorUserNotExecutor
	}
	if order.ExecutorID == order.CustomerID {
		return customErr.ErrorSameID
	}
	err = u.OrderRepo.UpdateExecutor(order, ctx)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (u *UseCase) DeleteExecutor(order models.Order, ctx context.Context) error {
	order.ExecutorID = 0
	err := u.OrderRepo.UpdateExecutor(order, ctx)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (u *UseCase) CloseOrder(orderID uint64, ctx context.Context) error {
	order, err := u.OrderRepo.FindByID(orderID, ctx)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	err = u.OrderRepo.DeleteOrder(orderID, ctx)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	_, err = u.OrderRepo.CreateArchive(*order, ctx)
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	return nil
}

func (u *UseCase) GetArchiveOrders(ctx context.Context) ([]models.Order, error) {
	orders, err := u.OrderRepo.GetArchiveOrders(ctx)
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

//TODO: вынести в отдеьлный модуль
func (u *UseCase) validateOrder(order *models.Order) error {
	err := validation.ValidateStruct(
		order,
		validation.Field(&order.OrderName, validation.Required, validation.Length(1, 300)),
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
	userR, err := u.UserRepo.GetUserById(context.Background(), &api.UserRequest{Id: order.CustomerID})
	if err != nil {
		return errors.Wrap(err, orderUseCaseError)
	}
	order.UserLogin = userR.GetLogin()
	order.UserImg = userR.GetImg()
	return nil
}
