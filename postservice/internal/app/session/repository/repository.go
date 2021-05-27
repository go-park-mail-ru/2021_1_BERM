package repository

import (
	"context"
	"post/api"
	"post/internal/app/models"
	"time"
)

type Repository struct {
	client api.SessionClient
}

func New(client api.SessionClient) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) Check(ctx context.Context, sessionID string) (*models.UserBasicInfo, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	u, err := r.client.Check(timeOutCtx, &api.SessionCheckRequest{
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
