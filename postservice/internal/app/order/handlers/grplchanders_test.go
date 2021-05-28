package order_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"post/api"
	"post/internal/app/models"
	ordHandlers "post/internal/app/order/handlers"
	"post/internal/app/order/mock"
	"post/pkg/metric"
	"testing"
)

func TestGRPCServer_GetOrderById(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderInfo := models.Order{
		ID:          1,
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
		UserLogin:   "astlok",
	}

	expectResponse := &api.OrderInfoResponse{
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
	}
	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := ordHandlers.NewGRPCServer(mockUseCase)

	mockUseCase.EXPECT().
		FindByID(uint64(1), context.Background()).
		Times(1).
		Return(&orderInfo, nil)

	response, err := handle.GetOrderById(context.Background(), &api.OrderRequest{Id: 1})
	require.NoError(t, err)
	require.Equal(t, expectResponse, response)
	metric.Destroy()
}

func TestGRPCServer_GetOrderByIdErr(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderInfo := models.Order{
		ID:          1,
		OrderName:   "Сверстать сайт",
		Category:    "Back",
		CustomerID:  1,
		Deadline:    1617004533,
		Budget:      1488,
		Description: "Pomogite sdelat API",
		UserLogin:   "astlok",
	}

	expectResponse := &api.OrderInfoResponse{
		OrderName:   "",
		Category:    "",
		CustomerID:  0,
		Deadline:    0,
		Budget:      0,
		Description: "",
	}
	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := ordHandlers.NewGRPCServer(mockUseCase)

	mockUseCase.EXPECT().
		FindByID(uint64(1), context.Background()).
		Times(1).
		Return(&orderInfo, errors.New("ERR"))

	response, err := handle.GetOrderById(context.Background(), &api.OrderRequest{Id: 1})
	require.Error(t, err)
	require.Equal(t, expectResponse, response)
	metric.Destroy()
}
