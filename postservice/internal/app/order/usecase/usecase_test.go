package order

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"post/api"
	"post/internal/app/models"
	"post/internal/app/order/mock"
	"testing"
)

func TestDeleteExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := models.Order{
		ExecutorID: 0,
	}
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockOrderRepo.EXPECT().UpdateExecutor(order, ctx).Times(1).Return(nil)

	mockUserRepo := mock.NewMockUserClient(ctrl)

	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	err := useCase.DeleteExecutor(order, ctx)

	require.NoError(t, err)
}

func TestDeleteExecutorErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := models.Order{
		ExecutorID: 0,
	}
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockOrderRepo.EXPECT().UpdateExecutor(order, ctx).Times(1).Return(errors.New("Db dead"))
	mockUserRepo := mock.NewMockUserClient(ctrl)

	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	err := useCase.DeleteExecutor(order, ctx)
	require.Error(t, err)
}

func TestCreateOrder (t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := models.Order{
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
	}
	expectOrder := &models.Order{
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
		UserImg: "kek",
		UserLogin: "Mem",
		ID: 1,
	}
	var id uint64
	id = 1
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)


	mockOrderRepo.EXPECT().
		Create(order, ctx).
		Times(1).
		Return(id, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respOrder, err := useCase.Create(order, ctx)

	require.Equal(t, respOrder, expectOrder)
	require.NoError(t, err)


	id = 0
	mockOrderRepo.EXPECT().
		Create(order, ctx).
		Times(1).
		Return(id, errors.New("DB error"))

	respOrder, err = useCase.Create(order, ctx)
	require.Error(t, err)


	id = 1
	mockOrderRepo.EXPECT().
		Create(order, ctx).
		Times(1).
		Return(id, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC Err"))

	respOrder, err = useCase.Create(order, ctx)
	require.Error(t, err)


	id = 1
	wrongOrder := models.Order{
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Category: "Back",
	}

	respOrder, err = useCase.Create(wrongOrder, ctx)
	require.Error(t, err)
}

func TestFindByID (t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := &models.Order{
		ID: 1,
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
	}
	expectOrder := &models.Order{
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
		UserImg: "kek",
		UserLogin: "Mem",
		ID: 1,
	}
	var id uint64
	id = 1
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)


	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(order, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respOrder, err := useCase.FindByID(id, ctx)

	require.Equal(t, respOrder, expectOrder)
	require.NoError(t, err)


	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(nil, nil)
	mockOrderRepo.EXPECT().
		FindArchiveByID(id, ctx).
		Times(1).
		Return(order, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respOrder, err = useCase.FindByID(id, ctx)

	require.Equal(t, respOrder, expectOrder)
	require.NoError(t, err)


	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(order, errors.New("DB err"))

	respOrder, err = useCase.FindByID(id, ctx)

	require.Error(t, err)

	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(nil, nil)
	mockOrderRepo.EXPECT().
		FindArchiveByID(id, ctx).
		Times(1).
		Return(nil, nil)

	respOrder, err = useCase.FindByID(id, ctx)

	require.NoError(t, err)



	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(nil, nil)
	mockOrderRepo.EXPECT().
		FindArchiveByID(id, ctx).
		Times(1).
		Return(order, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	respOrder, err = useCase.FindByID(id, ctx)

	require.Error(t, err)
}

func TestFindByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var id uint64
	id = 1
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	orders := []models.Order{
		{
			ID: 1,
			OrderName: "Keke",
			CustomerID: 1,
			Budget: 1488,
			Deadline: 22842212,
			Description: "Aue jizn voram",
			Category: "Back",
		},
	}

	expectOrders := []models.Order{
		{
			ID: 1,
			OrderName: "Keke",
			CustomerID: 1,
			Budget: 1488,
			Deadline: 22842212,
			Description: "Aue jizn voram",
			Category: "Back",
		},
	}

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(2).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)
	mockOrderRepo.EXPECT().
		FindByExecutorID(id, ctx).
		Times(1).
		Return(orders, nil)

	respOrders, err := useCase.FindByUserID(id, ctx)

	require.Equal(t, respOrders, expectOrders)
	require.NoError(t, err)


	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, errors.New("GRPC Err"))

	respOrders, err = useCase.FindByUserID(id, ctx)

	require.Error(t, err)


	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(2).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)
	mockOrderRepo.EXPECT().
		FindByCustomerID(id, ctx).
		Times(1).
		Return(orders, nil)

	respOrders, err = useCase.FindByUserID(id, ctx)

	require.Equal(t, respOrders, expectOrders)
	require.NoError(t, err)

	 expectEmptyOrders := []models.Order{}
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)
	mockOrderRepo.EXPECT().
		FindByCustomerID(id, ctx).
		Times(1).
		Return(nil, nil)

	respOrders, err = useCase.FindByUserID(id, ctx)

	require.Equal(t, respOrders, expectEmptyOrders)
	require.NoError(t, err)


	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)
	mockOrderRepo.EXPECT().
		FindByCustomerID(id, ctx).
		Times(1).
		Return(nil, errors.New("DB err"))

	respOrders, err = useCase.FindByUserID(id, ctx)

	require.Error(t, err)
}

func TestChangeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := &models.Order{
		ID: 1,
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
	}
	oldOrder := &models.Order{
		ID: 1,
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
	}
	expectOrder := &models.Order{
		OrderName: "Keke",
		CustomerID: 1,
		Budget: 1488,
		Deadline: 22842212,
		Description: "Aue jizn voram",
		Category: "Back",
		UserImg: "kek",
		UserLogin: "Mem",
		ID: 1,
	}
	var id uint64
	id = 1
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)


	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldOrder, nil)
	mockOrderRepo.EXPECT().
		Change(*order, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respOrder, err := useCase.ChangeOrder(*order, ctx)

	require.Equal(t, respOrder, *expectOrder)
	require.NoError(t, err)


	orderWithoutFields := &models.Order{
		ID: 1,
		CustomerID: 1,
	}
	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldOrder, nil)
	mockOrderRepo.EXPECT().
		Change(*order, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respOrder, err = useCase.ChangeOrder(*orderWithoutFields, ctx)

	require.Equal(t, respOrder, *expectOrder)
	require.NoError(t, err)



	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldOrder, errors.New("DB err"))

	respOrder, err = useCase.ChangeOrder(*orderWithoutFields, ctx)

	require.Error(t, err)
}