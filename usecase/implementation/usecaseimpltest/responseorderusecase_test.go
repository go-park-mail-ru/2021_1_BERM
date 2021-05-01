package usecaseimpl

import (
	"FL_2/model"
	"FL_2/store/mock"
	"FL_2/usecase/implementation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	correctResponseOrderModel = model.ResponseOrder{
		OrderID:   1,
		UserID:    1,
		Time:      123,
		Rate:      1,
		UserLogin: correctLogin,
		UserImg:   "",
	}
)

func TestResponseOrderCreate(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindSpecializesByUserID(uint64(1))
	responseOrderRepMock := mock.NewMockResponseOrderRepository(ctrl)
	responseOrderRepMock.EXPECT().Create(correctResponseOrderModel).Return(uint64(1), nil)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)
	mockStore.EXPECT().ResponseOrder().Times(1).Return(responseOrderRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.ResponseOrder().Create(correctResponseOrderModel)
	require.NoError(t, err)
}

func TestResponseOrderFindByVacancyID(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindSpecializesByUserID(uint64(1))
	responseOrderRepMock := mock.NewMockResponseOrderRepository(ctrl)
	responseOrderRepMock.EXPECT().FindByOrderID(uint64(1)).Return([]model.ResponseOrder{correctResponseOrderModel}, nil)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)
	mockStore.EXPECT().ResponseOrder().Times(1).Return(responseOrderRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.ResponseOrder().FindByVacancyID(1)
	require.NoError(t, err)
}

func TestResponseChange(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectOrderResponse := &model.ResponseOrder{}
	*newCorrectOrderResponse = correctResponseOrderModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindSpecializesByUserID(uint64(1))

	responseOrderRepMock := mock.NewMockResponseOrderRepository(ctrl)
	responseOrderRepMock.EXPECT().Change(correctResponseOrderModel).Return(newCorrectOrderResponse, nil)

	mockStore.EXPECT().User().Times(2).Return(userRepMock)
	mockStore.EXPECT().ResponseOrder().Times(1).Return(responseOrderRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.ResponseOrder().Change(correctResponseOrderModel)
	require.NoError(t, err)
}

func TestResponseDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	responseOrderRepMock := mock.NewMockResponseOrderRepository(ctrl)
	responseOrderRepMock.EXPECT().Delete(correctResponseOrderModel).Return(nil)

	mockStore.EXPECT().ResponseOrder().Times(1).Return(responseOrderRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	err := useCase.ResponseOrder().Delete(correctResponseOrderModel)
	require.NoError(t, err)
}
