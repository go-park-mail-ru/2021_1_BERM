package order

import (
	"context"
	"post/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock post/internal/app/order UseCase
type UseCase interface {
	Create(order models.Order, ctx context.Context) (*models.Order, error)
	FindByID(id uint64, ctx context.Context) (*models.Order, error)
	FindByUserID(userID uint64, ctx context.Context) ([]models.Order, error)
	ChangeOrder(order models.Order, ctx context.Context) (models.Order, error)
	DeleteOrder(id uint64, ctx context.Context) error
	GetActualOrders(ctx context.Context) ([]models.Order, uint64, error)
	SelectExecutor(order models.Order, ctx context.Context) error
	DeleteExecutor(order models.Order, ctx context.Context) error
	CloseOrder(orderID uint64, ctx context.Context) error
	GetArchiveOrders(userInfo models.UserBasicInfo, ctx context.Context) ([]models.Order, error)
	SearchOrders(keyword string, ctx context.Context) ([]models.Order, error)
}
