package handlers_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"user/api"
	"user/internal/app/models"
	userHandlers "user/internal/app/user/handlers"
	"user/internal/app/user/mock"
	"user/pkg/metric"
)

func TestGRPCServer_Login(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pass := "123456790dsxf"
	email := "dada@mail.ru"

	ctx := context.Background()
	mockUserUseCase := mock.NewMockUseCase(ctrl)
	mockUserUseCase.EXPECT().Verification(email, pass, ctx).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)

	handle := userHandlers.NewGRPCServer(mockUserUseCase)
	response, err := handle.AuthorizationUser(ctx, &api.AuthorizationUserRequest{
		Password: pass,
		Email:    email,
		ReqId:    1313213,
	})
	require.NoError(t, err)
	require.Equal(t, response.Id, uint64(1))
	metric.Destroy()
}

func TestGRPCServer_Registration(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pass := "123456790dsxf"
	email := "dada@mail.ru"

	ctx := context.Background()
	mockUserUseCase := mock.NewMockUseCase(ctrl)
	newUser := models.NewUser{
		Email:       email,
		Login:       "log",
		NameSurname: "log log",
		Password:    pass,
		About:       "asdsadsadsa",
	}
	mockUserUseCase.EXPECT().Create(newUser, ctx).Times(1).Return(&models.UserBasicInfo{
		ID:       1,
		Executor: true,
	}, nil)

	handle := userHandlers.NewGRPCServer(mockUserUseCase)
	response, err := handle.RegistrationUser(ctx, &api.NewUserRequest{
		About:       newUser.About,
		Email:       newUser.Email,
		Password:    pass,
		Login:       newUser.Login,
		NameSurname: newUser.NameSurname,
	})
	require.NoError(t, err)
	require.Equal(t, response.Id, uint64(1))
	metric.Destroy()
}
