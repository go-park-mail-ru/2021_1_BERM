package order

import (
	"context"
	"post/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock post/internal/app/order Repository
type Repository interface {
	Create(order models.Order, ctx context.Context) (uint64, error)
	Change(order models.Order, ctx context.Context) error
	DeleteOrder(id uint64, ctx context.Context) error
	FindByID(id uint64, ctx context.Context) (*models.Order, error)
	FindByExecutorID(executorID uint64, ctx context.Context) ([]models.Order, error)
	FindByCustomerID(customerID uint64, ctx context.Context) ([]models.Order, error)
	GetActualOrders(ctx context.Context) ([]models.Order, error)
	UpdateExecutor(order models.Order, ctx context.Context) error
	CreateArchive(order models.Order, ctx context.Context) error
	GetArchiveOrdersByExecutorID(executorID uint64, ctx context.Context) ([]models.Order, error)
	GetArchiveOrdersByCustomerID(customerID uint64, ctx context.Context) ([]models.Order, error)
	SearchOrders(keyword string, ctx context.Context) ([]models.Order, error)
	FindArchiveByID(id uint64, ctx context.Context) (*models.Order, error)
	SuggestOrderTitle(suggestWord string, ctx context.Context) ([]models.SuggestOrderTitle, error)
}
