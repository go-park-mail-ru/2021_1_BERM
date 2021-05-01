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
	correctVacancyModel = model.Vacancy{
		UserID:      1,
		Category:    "develop",
		VacancyName: ":ALALASDLASDDA:KM",
		Description: "SADJSAHJA:S",
		Salary:      5,
		Login:       correctLogin,
		Img:         "",
	}
)

func TestVacancyCreate(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	vacancyRepMock := mock.NewMockVacancyRepository(ctrl)
	vacancyRepMock.EXPECT().Create(correctVacancyModel).Return(uint64(1), nil)

	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)

	mockStore.EXPECT().Vacancy().Times(1).Return(vacancyRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Vacancy().Create(correctVacancyModel)
	require.NoError(t, err)
}

func TestVacancyFindByID(t *testing.T) {
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectVacancy := &model.Vacancy{}
	*newCorrectVacancy = correctVacancyModel
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStore(ctrl)
	vacancyRepMock := mock.NewMockVacancyRepository(ctrl)
	vacancyRepMock.EXPECT().FindByID(uint64(1)).Return(newCorrectVacancy, nil)

	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(uint64(1)).Return(newCorrectUser, nil)

	mockStore.EXPECT().Vacancy().Times(1).Return(vacancyRepMock)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)

	mockMediaStore := mock.NewMockMediaStore(ctrl)
	imageRepMock := mock.NewMockImageRepository(ctrl)
	imageRepMock.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(imageRepMock)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	_, err := useCase.Vacancy().FindByID(1)
	require.NoError(t, err)
}
