package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"user/internal/app/models"
	mock3 "user/internal/app/review/mock"
	mock2 "user/internal/app/specialize/mock"
	"user/internal/app/user/mock"
	"user/internal/app/user/tools/passwordencrypt"
	customError "user/pkg/error"
)

//Проверка созданияю юзера клиента
func TestCreateUserClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUser := &models.NewUser{
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		About:       "sdaasd sadasdAS DSdaS DAS",
	}
	ctx := context.Background()
	mockUserRepo := mock.NewMockRepository(ctrl)
	mockEncrypter := mock.NewMockPasswordEncrypter(ctrl)
	bcCall := mockEncrypter.EXPECT().BeforeCreate(newUser).Times(1)
	newUser.EncryptPassword = []byte{1, 2, 3, 4, 5}
	newUser.Executor = true
	bcCall.Return(newUser, nil)
	mockUserRepo.EXPECT().Create(*newUser, ctx).Times(1).Return(uint64(1), nil)
	useCase := UseCase{
		userRepository: mockUserRepo,
		encrypter: mockEncrypter,
	}
	userBasicInfo, err := useCase.Create(*newUser, ctx)
	require.NoError(t, err)
	require.Equal(t, userBasicInfo.ID, uint64(1))
	require.Equal(t, userBasicInfo.Executor, false)
}

//Проверка созданияю юзера клиента с невалидными данными
func TestCreateUserClientWithInvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUser := &models.NewUser{
		Email:       "abru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxA",
		About:       "sdaasd sadasdAS DSdaS DAS",
	}
	ctx := context.Background()

	useCase := UseCase{}
	_, err := useCase.Create(*newUser, ctx)
	require.Error(t, err)

}

//Проверка созданияю юзера исполнителя при отсутствии в базе указанных им специализаций
func TestCreateUserExecutorWithoutSpecInDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := make(pq.StringArray, 2)
	for index, _ := range spec {
		spec[index] = "123" + strconv.Itoa(index)
	}
	newUser := &models.NewUser{
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		About:       "sdaasd sadasdAS DSdaS DAS",
		Specializes: spec,
	}
	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().Create(newUser, ctx).Times(1).Return(uint64(1), nil)

	mockSpecializeRepo := mock2.NewMockRepository(ctrl)
	mockSpecializeRepo.EXPECT().FindByName(spec[0], ctx).Times(1).Return(uint64(0), customError.ErrorNoRows)
	mockSpecializeRepo.EXPECT().FindByName(spec[1], ctx).Times(1).Return(uint64(0), customError.ErrorNoRows)
	mockSpecializeRepo.EXPECT().Create(spec[0], ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeRepo.EXPECT().Create(spec[1], ctx).Times(1).Return(uint64(2), nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(2), uint64(1), ctx).Times(1).Return(nil)
	useCase := UseCase{
		userRepository:       mockUserRepo,
		specializeRepository: mockSpecializeRepo,
	}
	userBasicInfo, err := useCase.Create(*newUser, ctx)
	require.NoError(t, err)
	require.Equal(t, userBasicInfo.ID, uint64(1))
	require.Equal(t, userBasicInfo.Executor, true)
}

//Проверка созданияю юзера исполнителя при наличии в базе указанных им специализаций
func TestCreateUserExecutorWithSpecInDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := make(pq.StringArray, 2)
	for index, _ := range spec {
		spec[index] = "123" + strconv.Itoa(index)
	}
	newUser := &models.NewUser{
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		About:       "sdaasd sadasdAS DSdaS DAS",
		Specializes: spec,
	}
	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().Create(newUser, ctx).Times(1).Return(uint64(1), nil)

	mockSpecializeRepo := mock2.NewMockRepository(ctrl)
	mockSpecializeRepo.EXPECT().FindByName(spec[0], ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeRepo.EXPECT().FindByName(spec[1], ctx).Times(1).Return(uint64(2), nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(2), uint64(1), ctx).Times(1).Return(nil)
	useCase := UseCase{
		userRepository:       mockUserRepo,
		specializeRepository: mockSpecializeRepo,
	}
	userBasicInfo, err := useCase.Create(*newUser, ctx)
	require.NoError(t, err)
	require.Equal(t, userBasicInfo.ID, uint64(1))
	require.Equal(t, userBasicInfo.Executor, true)
}

//Тестирование поиска информации о юзере по id, когда юзер является исоплнителем
func TestGetByIDExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := make(pq.StringArray, 2)
	for index, _ := range spec {
		spec[index] = "123" + strconv.Itoa(index)
	}
	userInfo := &models.UserInfo{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Specializes: spec,
		Rating:      3,
		ReviewCount: 2,
		Executor:    true,
	}

	userReviewInfo := &models.UserReviewInfo{
		ReviewCount: 2,
		Rating:      4.2,
	}

	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().FindUserByID(uint64(1), ctx).Times(1).Return(userInfo, nil)

	mockSpecializeRepo := mock2.NewMockRepository(ctrl)
	mockSpecializeRepo.EXPECT().FindByUserID(uint64(1), ctx).Times(1).Return(spec, nil)

	mockReviewRepo := mock3.NewMockRepository(ctrl)
	mockReviewRepo.EXPECT().GetAvgScoreByUserId(uint64(1), ctx).Times(1).Return(userReviewInfo, nil)

	useCase := UseCase{
		userRepository:       mockUserRepo,
		specializeRepository: mockSpecializeRepo,
		reviewsRepository:    mockReviewRepo,
	}
	userInfo, err := useCase.GetById(1, ctx)
	require.NoError(t, err)
	require.Equal(t, userInfo.ID, uint64(1))
	require.Equal(t, userInfo.ReviewCount, userReviewInfo.ReviewCount)
	require.Equal(t, userInfo.Rating, userReviewInfo.Rating)

}

//Тестирование поиска информации о юзере по id, когда юзер является исоплнителем
func TestGetByIDClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userInfo := &models.UserInfo{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Rating:      3,
		ReviewCount: 2,
		Executor:    false,
	}

	userReviewInfo := &models.UserReviewInfo{
		ReviewCount: 2,
		Rating:      4.2,
	}

	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().FindUserByID(uint64(1), ctx).Times(1).Return(userInfo, nil)

	mockReviewRepo := mock3.NewMockRepository(ctrl)
	mockReviewRepo.EXPECT().GetAvgScoreByUserId(uint64(1), ctx).Times(1).Return(userReviewInfo, nil)

	useCase := UseCase{
		userRepository:    mockUserRepo,
		reviewsRepository: mockReviewRepo,
	}
	userInfo, err := useCase.GetById(1, ctx)
	require.NoError(t, err)
	require.Equal(t, userInfo.ID, uint64(1))
	require.Equal(t, userInfo.ReviewCount, userReviewInfo.ReviewCount)
	require.Equal(t, userInfo.Rating, userReviewInfo.Rating)

}

//Тестирование изменение юзера при наличии в базе необходимых специадизаций
func TestChangeUserWithSpecInDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := make(pq.StringArray, 2)
	for index, _ := range spec {
		spec[index] = "123" + strconv.Itoa(index)
	}
	changeUser := &models.ChangeUser{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		NewPassword: "adsadasda211@Sdas",
		About:       "sdaasd sadasdAS DSdaS DAS",
		Specializes: spec,
	}
	newUser := &models.NewUser{
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		About:       "sdaasd sadasdAS DSdaS DAS",
		Specializes: spec,
	}
	encrypter := &passwordencrypt.PasswordEncrypter{}
	var err error
	*newUser, err = encrypter.BeforeCreate(*newUser)
	if err != nil {
		return
	}
	userInfo := &models.UserInfo{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    newUser.EncryptPassword,
		Rating:      3,
		ReviewCount: 2,
		Executor:    false,
	}

	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().FindUserByID(changeUser.ID, ctx).Times(1).Return(userInfo, nil)
	mockUserRepo.EXPECT().Change(changeUser, ctx).Times(1).Return(nil)

	mockSpecializeRepo := mock2.NewMockRepository(ctrl)
	mockSpecializeRepo.EXPECT().FindByName(spec[0], ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeRepo.EXPECT().FindByName(spec[1], ctx).Times(1).Return(uint64(2), nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(2), uint64(1), ctx).Times(1).Return(nil)
	useCase := UseCase{
		userRepository:       mockUserRepo,
		specializeRepository: mockSpecializeRepo,
	}
	userBasicInfo, err := useCase.Change(*changeUser, ctx)
	require.NoError(t, err)
	require.Equal(t, userBasicInfo.ID, uint64(1))
	require.Equal(t, userBasicInfo.Executor, true)
}

//Тестирование изменение юзера при отсутствии в базе необходимых специадизаций
func TestChangeUserWithoutSpecInDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := make(pq.StringArray, 2)
	for index, _ := range spec {
		spec[index] = "123" + strconv.Itoa(index)
	}
	changeUser := &models.ChangeUser{
		ID:       1,
		Password: "zxcv12345@asd;A",

		Specializes: spec,
	}
	newUser := &models.NewUser{
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		About:       "sdaasd sadasdAS DSdaS DAS",
		Specializes: spec,
	}
	encrypter := &passwordencrypt.PasswordEncrypter{}
	var err error
	*newUser, err = encrypter.BeforeCreate(*newUser)
	if err != nil {
		return
	}
	spec1 := make(pq.StringArray, 1)
	spec1[0] = "ffsd"
	userInfo := &models.UserInfo{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    newUser.EncryptPassword,
		Specializes: spec1,
		Rating:      3,
		ReviewCount: 2,
		Executor:    false,
	}

	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().FindUserByID(changeUser.ID, ctx).Times(1).Return(userInfo, nil)
	mockUserRepo.EXPECT().Change(changeUser, ctx).Times(1).Return(nil)

	mockSpecializeRepo := mock2.NewMockRepository(ctrl)
	mockSpecializeRepo.EXPECT().FindByName(spec[0], ctx).Times(1).Return(uint64(0), customError.ErrorNoRows)
	mockSpecializeRepo.EXPECT().FindByName(spec[1], ctx).Times(1).Return(uint64(0), customError.ErrorNoRows)
	mockSpecializeRepo.EXPECT().FindByName(spec1[0], ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeRepo.EXPECT().Create(spec[0], ctx).Times(1).Return(uint64(2), nil)
	mockSpecializeRepo.EXPECT().Create(spec[1], ctx).Times(1).Return(uint64(3), nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(2), uint64(1), ctx).Times(1).Return(nil)
	mockSpecializeRepo.EXPECT().AssociateSpecializationWithUser(uint64(3), uint64(1), ctx).Times(1).Return(nil)
	useCase := UseCase{
		userRepository:       mockUserRepo,
		specializeRepository: mockSpecializeRepo,
	}
	userBasicInfo, err := useCase.Change(*changeUser, ctx)
	require.NoError(t, err)
	require.Equal(t, userBasicInfo.ID, uint64(1))
	require.Equal(t, userBasicInfo.Executor, true)
}

//Логин юзера с валидным поролем
func TestUserVerification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	spec := make(pq.StringArray, 2)
	for index, _ := range spec {
		spec[index] = "123" + strconv.Itoa(index)
	}
	newUser := &models.NewUser{
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    "zxcv12345@asd;A",
		About:       "sdaasd sadasdAS DSdaS DAS",
		Specializes: spec,
	}
	encrypter := &passwordencrypt.PasswordEncrypter{}
	var err error
	*newUser, err = encrypter.BeforeCreate(*newUser)
	if err != nil {
		return
	}
	userInfo := &models.UserInfo{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    newUser.EncryptPassword,
		Rating:      3,
		ReviewCount: 2,
		Executor:    false,
	}

	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().FindUserByEmail(newUser.Email, ctx).Times(1).Return(userInfo, nil)

	useCase := UseCase{
		userRepository: mockUserRepo,
	}
	userBasicInfo, err := useCase.Verification(newUser.Email, newUser.Password, ctx)
	require.NoError(t, err)
	require.Equal(t, userBasicInfo.ID, uint64(1))
	require.Equal(t, userBasicInfo.Executor, false)
}

//Логин юзера с невалидным поролем
func TestUserVerificationBadPass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userInfo := &models.UserInfo{
		ID:          1,
		Email:       "abc@mail.ru",
		Login:       "abcdefg",
		NameSurname: "abc bdf",
		Password:    []byte{1, 2, 4, 5, 6, 4, 2, 4, 6, 3, 45, 3, 3},
		Rating:      3,
		ReviewCount: 2,
		Executor:    false,
	}

	ctx := context.Background()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().FindUserByEmail("asdas@mail.ru", ctx).Times(1).Return(userInfo, nil)

	useCase := UseCase{
		userRepository: mockUserRepo,
	}
	_, err := useCase.Verification("asdas@mail.ru", "SAdadasdsda", ctx)
	require.Error(t, err)
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockRepository(ctrl)
	mockSpecializeRepo := mock2.NewMockRepository(ctrl)
	mockReviewRepo := mock3.NewMockRepository(ctrl)
	u := New(mockUserRepo, mockSpecializeRepo, mockReviewRepo)
	require.Equal(t, u.userRepository, mockUserRepo)
	require.Equal(t, u.specializeRepository, mockSpecializeRepo)
	require.Equal(t, u.reviewsRepository, mockReviewRepo)
}
