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

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := models.Order{
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	expectOrder := &models.Order{
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
		UserImg:     "kek",
		UserLogin:   "Mem",
		ID:          1,
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
		OrderName:  "Keke",
		CustomerID: 1,
		Budget:     1488,
		Deadline:   22842212,
		Category:   "Back",
	}

	respOrder, err = useCase.Create(wrongOrder, ctx)
	require.Error(t, err)
}

func TestFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := &models.Order{
		ID:          1,
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	expectOrder := &models.Order{
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
		UserImg:     "kek",
		UserLogin:   "Mem",
		ID:          1,
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
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
		},
	}

	expectOrders := []models.Order{
		{
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
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
		ID:          1,
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	oldOrder := &models.Order{
		ID:          1,
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	expectOrder := &models.Order{
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
		UserImg:     "kek",
		UserLogin:   "Mem",
		ID:          1,
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
		ID:         1,
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

	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldOrder, nil)
	mockOrderRepo.EXPECT().
		Change(*order, ctx).
		Times(1).
		Return(errors.New("DB err"))

	respOrder, err = useCase.ChangeOrder(*order, ctx)
	require.Error(t, err)

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
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	respOrder, err = useCase.ChangeOrder(*order, ctx)

	require.Error(t, err)
}

func TestDeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var id uint64
	id = 1
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	mockOrderRepo.EXPECT().
		DeleteOrder(id, ctx).
		Times(1).
		Return(nil)

	err := useCase.DeleteOrder(id, ctx)

	require.NoError(t, err)

	mockOrderRepo.EXPECT().
		DeleteOrder(id, ctx).
		Times(1).
		Return(errors.New("DB err"))

	err = useCase.DeleteOrder(id, ctx)

	require.Error(t, err)
}

func TestGetActualOrders(t *testing.T) {
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
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
		},
	}

	expectOrders := []models.Order{
		{
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
			UserLogin:   "Mem",
			UserImg:     "kek",
		},
	}

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)
	mockOrderRepo.EXPECT().
		GetActualOrders(ctx).
		Times(1).
		Return(orders, nil)

	respOrders, err := useCase.GetActualOrders(ctx)

	require.Equal(t, respOrders, expectOrders)
	require.NoError(t, err)

	mockOrderRepo.EXPECT().
		GetActualOrders(ctx).
		Times(1).
		Return(orders, errors.New("DB err"))

	respOrders, err = useCase.GetActualOrders(ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, errors.New("GRPC error"))
	mockOrderRepo.EXPECT().
		GetActualOrders(ctx).
		Times(1).
		Return(orders, nil)

	respOrders, err = useCase.GetActualOrders(ctx)

	require.Error(t, err)

	emptyOrders := []models.Order{}

	mockOrderRepo.EXPECT().
		GetActualOrders(ctx).
		Times(1).
		Return(nil, nil)

	respOrders, err = useCase.GetActualOrders(ctx)

	require.Equal(t, respOrders, emptyOrders)
	require.NoError(t, err)
}

func TestSelectExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	order := &models.Order{
		ID:          1,
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
		ExecutorID:  2,
	}

	mockOrderRepo.EXPECT().
		UpdateExecutor(*order, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)

	err := useCase.SelectExecutor(*order, ctx)

	require.NoError(t, err)

	errOrder := &models.Order{
		ID:          1,
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
		ExecutorID:  1,
	}

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: errOrder.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)

	err = useCase.SelectExecutor(*errOrder, ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: errOrder.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)

	err = useCase.SelectExecutor(*errOrder, ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, errors.New("GRPC err"))

	err = useCase.SelectExecutor(*order, ctx)

	require.Error(t, err)

	mockOrderRepo.EXPECT().
		UpdateExecutor(*order, ctx).
		Times(1).
		Return(errors.New("DB err"))
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)

	err = useCase.SelectExecutor(*order, ctx)

	require.Error(t, err)
}

func TestCloseOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	var id uint64
	id = 1

	order := &models.Order{
		ID:          1,
		OrderName:   "Keke",
		CustomerID:  1,
		Budget:      1488,
		Deadline:    22842212,
		Description: "Aue jizn voram",
		Category:    "Back",
		ExecutorID:  2,
	}

	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(order, nil)
	mockOrderRepo.EXPECT().
		DeleteOrder(id, ctx).
		Times(1).
		Return(nil)
	mockOrderRepo.EXPECT().
		CreateArchive(*order, ctx).
		Times(1).
		Return(nil)

	err := useCase.CloseOrder(id, ctx)

	require.NoError(t, err)

	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(order, nil)
	mockOrderRepo.EXPECT().
		DeleteOrder(id, ctx).
		Times(1).
		Return(nil)
	mockOrderRepo.EXPECT().
		CreateArchive(*order, ctx).
		Times(1).
		Return(errors.New("DB err"))

	err = useCase.CloseOrder(id, ctx)

	require.Error(t, err)

	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(order, nil)
	mockOrderRepo.EXPECT().
		DeleteOrder(id, ctx).
		Times(1).
		Return(errors.New("DB err"))

	err = useCase.CloseOrder(id, ctx)

	require.Error(t, err)

	mockOrderRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(order, errors.New("DB err"))

	err = useCase.CloseOrder(id, ctx)

	require.Error(t, err)
}

func TestGetArchiveOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	var id uint64
	id = 1

	order := []models.Order{
		{
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
			ExecutorID:  2,
		},
	}

	expectOrders := []models.Order{
		{
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
			ExecutorID:  2,
			UserImg:     "kek",
			UserLogin:   "Mem",
		},
	}

	userInfo := models.UserBasicInfo{ID: id, Executor: false}

	mockOrderRepo.EXPECT().
		GetArchiveOrdersByCustomerID(id, ctx).
		Times(1).
		Return(order, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	resOrders, err := useCase.GetArchiveOrders(userInfo, ctx)

	require.Equal(t, expectOrders, resOrders)
	require.NoError(t, err)

	id = 2
	userInfo = models.UserBasicInfo{ID: id, Executor: true}

	mockOrderRepo.EXPECT().
		GetArchiveOrdersByExecutorID(id, ctx).
		Times(1).
		Return(order, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	resOrders, err = useCase.GetArchiveOrders(userInfo, ctx)

	require.Equal(t, expectOrders, resOrders)
	require.NoError(t, err)

	mockOrderRepo.EXPECT().
		GetArchiveOrdersByExecutorID(id, ctx).
		Times(1).
		Return(order, errors.New("DB err"))

	resOrders, err = useCase.GetArchiveOrders(userInfo, ctx)

	require.Error(t, err)

	mockOrderRepo.EXPECT().
		GetArchiveOrdersByExecutorID(id, ctx).
		Times(1).
		Return(order, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: order[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	resOrders, err = useCase.GetArchiveOrders(userInfo, ctx)
	require.Error(t, err)

	mockOrderRepo.EXPECT().
		GetArchiveOrdersByExecutorID(id, ctx).
		Times(1).
		Return(nil, nil)

	emptyOrder := []models.Order{}
	resOrders, err = useCase.GetArchiveOrders(userInfo, ctx)
	require.Equal(t, emptyOrder, resOrders)
	require.NoError(t, err)
}

func TestSearchOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockOrderRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := NewUseCase(mockOrderRepo, mockUserRepo)

	keyword := "Aue"
	orders := []models.Order{
		{
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
			ExecutorID:  2,
		},
	}

	expectOrders := []models.Order{
		{
			ID:          1,
			OrderName:   "Keke",
			CustomerID:  1,
			Budget:      1488,
			Deadline:    22842212,
			Description: "Aue jizn voram",
			Category:    "Back",
			ExecutorID:  2,
			UserImg:     "kek",
			UserLogin:   "Mem",
		},
	}

	mockOrderRepo.EXPECT().
		SearchOrders(keyword, ctx).
		Times(1).
		Return(orders, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: orders[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	resOrders, err := useCase.SearchOrders(keyword, ctx)

	require.Equal(t, expectOrders, resOrders)
	require.NoError(t, err)

	mockOrderRepo.EXPECT().
		SearchOrders(keyword, ctx).
		Times(1).
		Return(nil, nil)

	resOrders, err = useCase.SearchOrders(keyword, ctx)

	emptyOrders := []models.Order{}
	require.Equal(t, emptyOrders, resOrders)
	require.NoError(t, err)

	mockOrderRepo.EXPECT().
		SearchOrders(keyword, ctx).
		Times(1).
		Return(orders, errors.New("DB error"))

	resOrders, err = useCase.SearchOrders(keyword, ctx)

	require.Error(t, err)

	mockOrderRepo.EXPECT().
		SearchOrders(keyword, ctx).
		Times(1).
		Return(orders, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: orders[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	resOrders, err = useCase.SearchOrders(keyword, ctx)

	require.Error(t, err)
}
