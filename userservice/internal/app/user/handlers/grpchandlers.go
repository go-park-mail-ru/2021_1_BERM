package handlers

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user/api"
	"user/internal/app/models"
	"user/internal/app/user"
	"user/pkg/error/errortools"
)

type GRPCServer struct {
	api.UnimplementedUserServer
	userUseCase user.UseCase
}

func NewGRPCServer(userUseCase user.UseCase) *GRPCServer {
	return &GRPCServer{
		userUseCase: userUseCase,
	}
}

func (s *GRPCServer) RegistrationUser(ctx context.Context, in *api.NewUserRequest) (*api.UserResponse, error) {
	u := &models.NewUser{
		Email:       in.GetEmail(),
		Login:       in.GetLogin(),
		NameSurname: in.GetNameSurname(),
		Password:    in.GetPassword(),
		About:       in.GetAbout(),
		Specializes: in.GetSpecializes(),
	}

	answer, err := s.userUseCase.Create(*u, ctx)
	if err != nil {
		errData, codeUint32 := errortools.ErrorHandle(err)
		code := codes.Code(codeUint32)
		b, jsonErr := json.Marshal(errData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &api.UserResponse{
			Id:       0,
			Executor: false,
		}, status.Error(code, string(b))
	}

	return &api.UserResponse{
		Id:       answer.ID,
		Executor: answer.Executor,
	}, nil
}

func (s *GRPCServer) AuthorizationUser(ctx context.Context, in *api.AuthorizationUserRequest) (*api.UserResponse, error) {
	answer, err := s.userUseCase.Verification(in.GetEmail(), in.GetPassword(), ctx)
	if err != nil {
		errData, codeUint32 := errortools.ErrorHandle(err)
		code := codes.Code(codeUint32)
		serializeErrData, jsonErr := json.Marshal(errData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &api.UserResponse{
			Id:       0,
			Executor: false,
		}, status.Error(code, string(serializeErrData))
	}
	return &api.UserResponse{
		Id:       answer.ID,
		Executor: answer.Executor,
	}, nil
}

func (s *GRPCServer) GetUserById(ctx context.Context, in *api.UserRequest) (*api.UserInfoResponse, error) {
	userInfo, err := s.userUseCase.GetById(in.GetId(), ctx)
	if err != nil {
		errData, codeUint32 := errortools.ErrorHandle(err)
		code := codes.Code(codeUint32)
		serializeErrData, jsonErr := json.Marshal(errData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &api.UserInfoResponse{
			Email:       "",
			Login:       "",
			NameSurname: "",
			About:       "",
			Specializes: nil,
			Executor:    false,
			Img:         "",
			Rating:      0,
		}, status.Error(code, string(serializeErrData))
	}
	return &api.UserInfoResponse{
		Email:       userInfo.Email,
		Login:       userInfo.Login,
		NameSurname: userInfo.NameSurname,
		About:       userInfo.About,
		Specializes: userInfo.Specializes,
		Executor:    userInfo.Executor,
		Img:         userInfo.Img,
		Rating:      int32(userInfo.Rating),
	}, nil
}

func (s *GRPCServer) GetSpecializeByUserId(ctx context.Context, in *api.UserRequest) (*api.GetUserSpecializeResponse, error) {

	return nil, nil
}

func (s *GRPCServer) SetImgUrl(ctx context.Context, in *api.SetImgUrlRequest) (*api.SetImgUrlResponse, error) {
	err := s.userUseCase.SetImg(in.GetId(), in.GetImgIrl(), ctx)
	if err != nil {
		errData, codeUint32 := errortools.ErrorHandle(err)
		code := codes.Code(codeUint32)
		serializeErrData, jsonErr := json.Marshal(errData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &api.SetImgUrlResponse{Successfully: false}, status.Error(code, string(serializeErrData))
	}
	return &api.SetImgUrlResponse{Successfully: true}, nil
}
