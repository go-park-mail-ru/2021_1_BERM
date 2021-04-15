package usecaseimpl

import (
	"FL_2/model"
	"FL_2/store/mock"
	"FL_2/usecase/implementation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMediaUseCaseSetImage(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().ChangeUser(newCorrectUser).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindUserByID(uint64(0)).Return(newCorrectUser, nil)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().SetImage("0", []byte{}).Return("", nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Media().SetImage(newCorrectUser,[]byte{})
	require.NoError(t, err)
}
