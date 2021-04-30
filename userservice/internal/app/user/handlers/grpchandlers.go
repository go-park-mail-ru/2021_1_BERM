package handlers

import (
	"context"
	"user/api"
	"user/internal/app/models"
	"user/internal/app/user/usecase"
)

type GRPCServer struct {
	api.UnimplementedUserServer
	userUseCase usecase.UseCase
}

func NewGRPCServer(userUseCase usecase.UseCase) *GRPCServer {
	return &GRPCServer{
		userUseCase: userUseCase,
	}
}

func (s *GRPCServer) RegistrationUser(ctx context.Context, in *api.NewUserRequest) (*api.UserResponse, error) {
	u := models.NewUser{
		Email:       in.GetEmail(),
		Login:       in.GetLogin(),
		NameSurname: in.GetNameSurname(),
		Password:    in.GetPassword(),
		About:       in.GetAbout(),
		Specializes: in.GetSpecializes(),
	}

	answer, err := s.userUseCase.Create(u, ctx)
	if err != nil {
		return nil, err
	}

	return &api.UserResponse{
		Id:       answer["id"].(uint64),
		Executor: answer["executor"].(bool),
	}, nil
}

func (s *GRPCServer) AuthorizationUser(ctx context.Context, in *api.AuthorizationUserRequest) (*api.UserResponse, error) {
	answer, err := s.userUseCase.Verification(in.GetEmail(), in.GetPassword(), ctx)
	if err != nil {
		return nil, err
	}
	return &api.UserResponse{
		Id:       answer["id"].(uint64),
		Executor: answer["executor"].(bool),
	}, nil
}

func (s *GRPCServer) GetUserById(ctx context.Context, in *api.UserRequest) (*api.UserInfoResponse, error) {
	userInfo, err := s.userUseCase.GetById(in.GetId(), ctx)
	if err != nil {
		return nil, err
	}
	return &api.UserInfoResponse{
		Email:       userInfo.Email,
		Login:       userInfo.Login,
		NameSurname: userInfo.NameSurname,
		About:       userInfo.About,
		Specializes: userInfo.Specializes,
		Executor:    userInfo.Executor,
		Img:         userInfo.Img,
		Rating:      userInfo.Rating,
	}, nil
}

func (s *GRPCServer) GetSpecializeByUserId(ctx context.Context, in *api.UserRequest) (*api.GetUserSpecializeResponse, error) {

	return nil, nil
}

func (s *GRPCServer) SetImgUrl(ctx context.Context, in *api.SetImgUrlRequest) (*api.SetImgUrlResponse, error) {
	err := s.userUseCase.SetImg(in.GetId(), in.GetImgIrl(), ctx)
	if err != nil {
		return &api.SetImgUrlResponse{Successfully: false}, err
	}
	return &api.SetImgUrlResponse{Successfully: true}, nil
}
