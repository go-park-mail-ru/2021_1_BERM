package grpcrepository

import (
	"context"
	"user/api"
	"user/internal/app/models"
)

type Repository struct {
	client api.OrderClient
}

func New(client api.OrderClient) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) GetByID(ID uint64, ctx context.Context)  (*models.OrderInfo, error){
	o, err := r.client.GetOrderById(ctx, &api.OrderRequest{
		Id: ID,
	})
	if err != nil {
		return nil, err
	}
	return &models.OrderInfo{
		OrderName: o.OrderName,
		CustomerId: o.CustomerID,
		ExecutorId: o.ExecutorID,
		Budget: o.Budget,
		Description: o.Description,
		DeadLine: o.Deadline,
		Category: o.Category,
	}, nil
}

