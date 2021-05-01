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
	correctResponseVacancyModel = model.ResponseVacancy{
		VacancyID: 1,
		UserID:    1,
		Time:      12131,
		UserLogin: correctLogin,
		Rate:      1,
		UserImg:   "",
	}
)

func TestResponseVacancyCreate(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)

	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindSpecializesByUserID(uint64(1))

	responseVacancyRepMock := mock.NewMockResponseVacancyRepository(ctrl)
	responseVacancyRepMock.EXPECT().Create(correctResponseVacancyModel).Return(uint64(1), nil)

	mockStore.EXPECT().User().Times(2).Return(userRepMock)
	mockStore.EXPECT().ResponseVacancy().Times(1).Return(responseVacancyRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.ResponseVacancy().Create(correctResponseVacancyModel)
	require.NoError(t, err)
}

func TestResponseFindByVacancyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)

	responseVacancyRepMock := mock.NewMockResponseVacancyRepository(ctrl)
	responseVacancyRepMock.EXPECT().FindByVacancyID(uint64(1)).Return([]model.ResponseVacancy{correctResponseVacancyModel}, nil)

	mockStore.EXPECT().ResponseVacancy().Times(1).Return(responseVacancyRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.ResponseVacancy().FindByVacancyID(1)
	require.NoError(t, err)
}

func TestResponseVacancyChange(t *testing.T) {
	newRecponseVacancy := &model.ResponseVacancy{}
	*newRecponseVacancy = correctResponseVacancyModel

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)

	responseVacancyRepMock := mock.NewMockResponseVacancyRepository(ctrl)
	responseVacancyRepMock.EXPECT().Change(correctResponseVacancyModel).Return(newRecponseVacancy, nil)
	mockStore.EXPECT().ResponseVacancy().Times(1).Return(responseVacancyRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.ResponseVacancy().Change(correctResponseVacancyModel)
	require.NoError(t, err)
}

func TestResponseVacancyDelete(t *testing.T) {
	newRecponseVacancy := &model.ResponseVacancy{}
	*newRecponseVacancy = correctResponseVacancyModel

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)

	responseVacancyRepMock := mock.NewMockResponseVacancyRepository(ctrl)
	responseVacancyRepMock.EXPECT().Delete(correctResponseVacancyModel).Return(nil)
	mockStore.EXPECT().ResponseVacancy().Times(1).Return(responseVacancyRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	err := useCase.ResponseVacancy().Delete(correctResponseVacancyModel)
	require.NoError(t, err)
}
