package handlers

import (
	"authorizationservice/api"
	session2 "authorizationservice/internal/app/session"
	"context"
)

type GRPCServer struct {
	api.UnimplementedSessionServer
	sessionUseCase session2.UseCase
}

func NewGRPCServer(sessionUseCase session2.UseCase) *GRPCServer {
	return &GRPCServer{
		sessionUseCase: sessionUseCase,
	}
}

func (s *GRPCServer) Check(ctx context.Context, in *api.SessionCheckRequest) (*api.SessionCheckResponse, error) {
	session, err := s.sessionUseCase.Get(in.GetSessionId(), ctx)
	if err != nil {

		return nil, err
	}

	return &api.SessionCheckResponse{
		ID:       session.UserId,
		Executor: session.Executor,
	}, nil
}
