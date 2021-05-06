package repository

import (
	"authorizationservice/api"
	models2 "authorizationservice/internal/app/models"
	"context"
)

type Repository struct {
	client api.UserClient
}

func New(client api.UserClient) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) Create(newUser models2.NewUser, ctx context.Context) (*models2.UserBasicInfo, error) {
	//timeOutCtx, cancel := context.WithTimeout(ctx, time.Second)
	//defer cancel()
	userResponse, err := r.client.RegistrationUser(context.Background(), &api.NewUserRequest{
		Email:       newUser.Email,
		Login:       newUser.Login,
		NameSurname: newUser.NameSurname,
		Password:    newUser.Password,
		About:       newUser.About,
		Specializes: newUser.Specializes,
	})
	if err != nil {
		return nil, err
	}
	return &models2.UserBasicInfo{
		ID:       userResponse.GetId(),
		Executor: userResponse.GetExecutor(),
	}, err
}

func (r *Repository) Authentication(email string, password string, ctx context.Context) (*models2.UserBasicInfo, error) {
	//timeOutCtx, cancel := context.WithTimeout(ctx, time.Second)
	//defer cancel()
	userResponse, err := r.client.AuthorizationUser(context.Background(), &api.AuthorizationUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &models2.UserBasicInfo{
		ID:       userResponse.GetId(),
		Executor: userResponse.GetExecutor(),
	}, err
}
