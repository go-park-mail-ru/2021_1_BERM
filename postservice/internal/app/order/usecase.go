package order

import (
	"context"
	"post/internal/app/models"
)

type UseCase interface {
	Create(order models.Order, ctx context.Context) (*models.Order, error)
	FindByID(id uint64, ctx context.Context) (*models.Order, error)
	FindByUserID(userID uint64, ctx context.Context) ([]models.Order, error)
	ChangeOrder(order models.Order, ctx context.Context) (models.Order, error)
	DeleteOrder(id uint64, ctx context.Context) error
	GetActualOrders(ctx context.Context) ([]models.Order, error)
	SelectExecutor(order models.Order, ctx context.Context) error
	DeleteExecutor(order models.Order, ctx context.Context) error
	CloseOrder(orderID uint64, ctx context.Context) error
	GetArchiveOrders(userInfo models.UserBasicInfo, ctx context.Context) ([]models.Order, error)
	SearchOrders(keyword string, ctx context.Context) ([]models.Order, error)
}
