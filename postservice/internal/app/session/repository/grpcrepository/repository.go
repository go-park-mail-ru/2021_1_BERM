package grpcrepository

import (
	"context"
	"post/api"
	"post/internal/app/models"
)

type Repository struct {
	client api.SessionClient
}

func New(client api.SessionClient) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) Check(sessionID string, ctx context.Context) (*models.UserBasicInfo, error) {
	u, err := r.client.Check(ctx, &api.SessionCheckRequest{
		SessionId: sessionID,
	})
	if err != nil {
		return nil, err
	}
	return &models.UserBasicInfo{
		ID:       u.GetID(),
		Executor: u.Executor,
	}, err
}
