package usecaseimpl_test

import (
	"FL_2/model"
	"FL_2/store/mock"
	"FL_2/usecase/implementation"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"math/rand"
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
)

func TestCreateWithCorrectData(t *testing.T) {
	userId := uint64(1)
	specId := uint64(1)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().AddUser(newCorrectUser).Return(userId, nil)
	userRepMock.EXPECT().FindSpecializeByName(correctUser.Specializes[0]).Return(model.Specialize{
		ID:   specId,
		Name: correctUser.Specializes[0],
	}, nil)
	userRepMock.EXPECT().AddUserSpec(userId, specId).Return(nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(3).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.Create(newCorrectUser)
	require.NoError(t, err)
}

func TestCreateWithIncorrectPass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectUser.Password = "123"
	mockStore := mock.NewMockStore(ctrl)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.Create(newCorrectUser)
	require.Error(t, err)
}

func TestCreateWithMissingSpecialization(t *testing.T) {
	userId := uint64(1)
	specId := uint64(1)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().AddUser(newCorrectUser).Return(userId, nil)
	userRepMock.EXPECT().FindSpecializeByName(correctUser.Specializes[0]).Return(model.Specialize{
		ID:   0,
		Name: "",
	}, nil)
	userRepMock.EXPECT().AddSpec(correctUser.Specializes[0]).Return(specId, nil)
	userRepMock.EXPECT().AddUserSpec(userId, specId).Return(nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(4).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.Create(newCorrectUser)
	require.NoError(t, err)
}

func TestCreateWithoutSpec(t *testing.T) {
	userId := uint64(1)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectUser.Specializes = []string{}
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().AddUser(newCorrectUser).Return(userId, nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.Create(newCorrectUser)
	require.NoError(t, err)
}

func TestCreateWithDeadDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectUser.Specializes = []string{}
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().AddUser(newCorrectUser).Return(uint64(0), errors.New("Db dead"))
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.Create(newCorrectUser)
	require.Error(t, err)
}

func TestUserVerificationWithCorrectData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectUser.Specializes = nil
	newCorrectUser.Executor = false
	userRepMock := mock.NewMockUserRepository(ctrl)
	resUser := *newCorrectUser;
	salt := make([]byte, 8)
	rand.Read(salt)
	resUser.EncryptPassword = implementation.HashPass(salt, newCorrectUser.Password)
	userRepMock.EXPECT().FindUserByEmail(newCorrectUser.Email).Return(&resUser, nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	mockImageRep := mock.NewMockImageRepository(ctrl)
	mockImageRep.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(mockImageRep)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	u, err := userUseCase.UserVerification(newCorrectUser.Email, newCorrectUser.Password)
	require.NoError(t, err)
	require.Equal(t, *u, resUser)
}

func TestUserVerificationWithBadDb(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	userRepMock := mock.NewMockUserRepository(ctrl)
	resUser := *newCorrectUser;
	salt := make([]byte, 8)
	rand.Read(salt)
	resUser.EncryptPassword = implementation.HashPass(salt, newCorrectUser.Password)
	userRepMock.EXPECT().FindUserByEmail(newCorrectUser.Email).Return(nil,errors.New("err"))
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	u, err := userUseCase.UserVerification(newCorrectUser.Email, newCorrectUser.Password)
	require.Error(t, err)
	require.Nil(t, u)
}

func TestUserFindByIdWithCorrectDataNotExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	newCorrectUser.Specializes = nil
	newCorrectUser.Executor = false
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(newCorrectUser.ID).Return(newCorrectUser, nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(1).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	mockImageRep := mock.NewMockImageRepository(ctrl)
	mockImageRep.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(mockImageRep)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	_, err := userUseCase.FindByID(newCorrectUser.ID)
	require.NoError(t, err)
}

func TestUserFindByIdWithCorrectDataExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(newCorrectUser.ID).Return(newCorrectUser, nil)
	userRepMock.EXPECT().FindSpecializesByUserID(newCorrectUser.ID).Return([]string{}, nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	mockImageRep := mock.NewMockImageRepository(ctrl)
	mockImageRep.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(mockImageRep)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	_, err := userUseCase.FindByID(newCorrectUser.ID)
	require.NoError(t, err)
}

func TestUserChangeUserWithCorrectData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	resUser := &model.User{}
	*resUser = *newCorrectUser
	salt := make([]byte, 8)
	rand.Read(salt)
	resUser.EncryptPassword = implementation.HashPass(salt, newCorrectUser.Password)
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindUserByID(newCorrectUser.ID).Return(resUser, nil)
	resUser.Specializes = nil
	userRepMock.EXPECT().ChangeUser(newCorrectUser).Return(resUser, nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)
	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)
	mockImageRep := mock.NewMockImageRepository(ctrl)
	mockImageRep.EXPECT().GetImage("").Return(nil, nil)
	mockMediaStore.EXPECT().Image().Return(mockImageRep)
	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	_, err := userUseCase.ChangeUser(newCorrectUser)
	require.NoError(t, err)
}

func TestUserAssSpecializeWithCorrectData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindSpecializeByName(newCorrectUser.Specializes[0]).Return(model.Specialize{
		ID: 0,
		Name: "",
	}, nil)
	userRepMock.EXPECT().AddSpec(newCorrectUser.Specializes[0]).Return(uint64(1), nil)
	userRepMock.EXPECT().IsUserHaveSpec(uint64(1), uint64(1)).Return(false, nil)
	userRepMock.EXPECT().AddUserSpec(uint64(1), uint64(1)).Return(nil)
	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(4).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.AddSpecialize(newCorrectUser.Specializes[0], 1)
	require.NoError(t, err)
}

func TestUserpecDelWithCorrectData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newCorrectUser := &model.User{}
	*newCorrectUser = correctUser
	userRepMock := mock.NewMockUserRepository(ctrl)
	userRepMock.EXPECT().FindSpecializeByName(newCorrectUser.Specializes[0]).Return(model.Specialize{
		ID: 1,
		Name: newCorrectUser.Specializes[0] ,
	}, nil)
	userRepMock.EXPECT().DelSpecialize(uint64(1), uint64(1)).Return(nil)

	mockStore := mock.NewMockStore(ctrl)
	mockStore.EXPECT().User().Times(2).Return(userRepMock)

	mockCache := mock.NewMockCaсhe(ctrl)
	mockMediaStore := mock.NewMockMediaStore(ctrl)

	useCase := implementation.New(mockStore, mockCache, mockMediaStore)
	userUseCase := useCase.User()
	err := userUseCase.DelSpecialize(newCorrectUser.Specializes[0], 1)
	require.NoError(t, err)
}

//func (u *UserUseCase) DelSpecialize(specName string, userID uint64) error {
//	specialize, err := u.store.User().FindSpecializeByName(specName)
//	if err != nil {
//		return errors.Wrap(err, userUseCaseError)
//	}
//	err = u.store.User().DelSpecialize(specialize.ID, userID)
//
//	if err != nil {
//		return errors.Wrap(err, userUseCaseError)
//	}
//
//	return nil
//}
