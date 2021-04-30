package handlers

import (
	"authorizationservice/api"
	"authorizationservice/internal/session/usecase"
	"context"
)

type GRPCServer struct {
	api.UnimplementedSessionServer
	sessionUseCase usecase.UseCase
}

func NewGRPCServer(sessionUseCase usecase.UseCase) *GRPCServer {
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
