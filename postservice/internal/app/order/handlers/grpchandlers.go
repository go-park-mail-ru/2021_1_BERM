package order

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"post/api"
	"post/internal/app/order"
	"post/pkg/error/errortools"
)

type GRPCServer struct {
	api.UnimplementedOrderServer
	orderUseCase order.UseCase
}

func NewGRPCServer(orderUseCase order.UseCase) *GRPCServer {
	return &GRPCServer{
		orderUseCase: orderUseCase,
	}
}

func (s *GRPCServer) GetOrderById(ctx context.Context, in *api.OrderRequest) (*api.OrderInfoResponse, error) {
	orderInfo, err := s.orderUseCase.FindByID(in.GetId(), ctx)
	if err != nil {
		errData, codeUint32 := errortools.ErrorHandle(err)
		code := codes.Code(codeUint32)
		serializeErrData, jsonErr := json.Marshal(errData)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &api.OrderInfoResponse{
			OrderName:   "",
			CustomerID:  0,
			ExecutorID:  0,
			Budget:      0,
			Deadline:    0,
			Description: "",
			Category:    "",
		}, status.Error(code, string(serializeErrData))
	}
	return &api.OrderInfoResponse{
		OrderName:   orderInfo.OrderName,
		CustomerID:  orderInfo.CustomerID,
		ExecutorID:  orderInfo.ExecutorID,
		Budget:      orderInfo.Budget,
		Deadline:    orderInfo.Deadline,
		Description: orderInfo.Description,
		Category:    orderInfo.Category,
	}, nil
}
