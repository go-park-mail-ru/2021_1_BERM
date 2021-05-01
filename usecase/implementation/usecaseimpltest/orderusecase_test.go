package usecaseimpl

import (
	"FL_2/model"
	"FL_2/store/mock"
	"FL_2/usecase/implementation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	correctEmail    = "Email@gmail.com"
	correctPassword = "Zxcv1234"
	correctLogin    = "ValintinGovna"
	correctName     = "Valintin Valeninivich Gey"
	correctAbout    = "I am like it pool"
)

var (
	correctUser = model.User{
		Email:       correctEmail,
		Password:    correctPassword,
		Login:       correctLogin,
		NameSurname: correctName,
		Executor:    true,
		About:       correctAbout,
		Specializes: []string{
			"developer",
		},
	}

	correctOrderModel = model.Order{
		ID:          0,
		OrderName:   "123sdadssadqw",
		CustomerID:  1,
		ExecutorID:  2,
		Budget:      10,
		Deadline:    2002234,
		Description: "sdsdaDsDASADsADDD",
		Category:    "develop",
		Login:       "Valia",
		Img:         "",
	}
)

func TestOrderCreate(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().Create(correctOrderModel).Return(uint64(1), nil)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Order().Create(correctOrderModel)
	require.NoError(t, err)
}

func TestOrderFindById(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectOrder := &model.Order{}
	*newCorrectOrder = correctOrderModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().FindByID(uint64(1)).Return(newCorrectOrder, nil)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Order().FindByID(1)
	require.NoError(t, err)
}

func TestOrderFindByExecutorId(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectOrder := &model.Order{}
	*newCorrectOrder = correctOrderModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().FindByExecutorID(uint64(1)).Return([]model.Order{*newCorrectOrder}, nil)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Order().FindByExecutorID(1)
	require.NoError(t, err)
}

func TestOrderFindByCustomerID(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectOrder := &model.Order{}
	*newCorrectOrder = correctOrderModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().FindByCustomerID(uint64(1)).Return([]model.Order{*newCorrectOrder}, nil)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Order().FindByCustomerID(1)
	require.NoError(t, err)
}

func TestOrderGetActualOrders(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectOrder := &model.Order{}
	*newCorrectOrder = correctOrderModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().GetActualOrders().Return([]model.Order{*newCorrectOrder}, nil)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Order().GetActualOrders()
	require.NoError(t, err)
}

func TestOrderSelectExecutor(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectUser.ID = 2
	newCorrectOrder := &model.Order{}
	*newCorrectOrder = correctOrderModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().UpdateExecutor(*newCorrectOrder).Return(nil)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(2)).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindSpecializesByUserID(uint64(2)).Return([]string{"developer"}, nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	err := useCase.Order().SelectExecutor(*newCorrectOrder)
	require.NoError(t, err)
}

func TestOrderDelExecutor(t *testing.T) {
	newCorrectOrder := &model.Order{}
	*newCorrectOrder = correctOrderModel
	newCorrectOrder.ExecutorID = 0
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	orderRepMock := mock.NewMockOrderRepository(ctrl)
	orderRepMock.EXPECT().UpdateExecutor(*newCorrectOrder).Return(nil)
	mockStore.EXPECT().Order().Times(1).Return(orderRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	err := useCase.Order().DeleteExecutor(*newCorrectOrder)
	require.NoError(t, err)
}
